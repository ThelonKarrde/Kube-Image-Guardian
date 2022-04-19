package validation

import (
	"encoding/json"
	"fmt"
	"github.com/ThelonKarrde/Kube-Image-Guardian/internal/config"
	"github.com/gofiber/fiber/v2"
	"k8s.io/api/admission/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"strings"
)

type Validator struct {
	allowConfig config.AllowConfig
	infoLog     *log.Logger
	errorLog    *log.Logger
}

func (v *Validator) New(configPath string, infoLog *log.Logger, errorLog *log.Logger) {
	v.readConfig(configPath)
	v.infoLog = infoLog
	v.errorLog = errorLog
}

func (v *Validator) readConfig(path string) {
	v.allowConfig.ReadConfig(path)
}

func (v *Validator) Validate(c *fiber.Ctx) error {
	log.Print(v.allowConfig)
	var admissionReview v1beta1.AdmissionReview
	var err error
	err = json.Unmarshal(c.Body(), &admissionReview)
	if err != nil {
		v.errorLog.Println("Error decoding request data")
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(&fiber.Map{
			"error": "decoding request data",
		})
	}
	pod := v1.Pod{}
	err = json.Unmarshal(admissionReview.Request.Object.Raw, &pod)
	if err != nil {
		v.errorLog.Println("Error unmarshalling request data")
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(&fiber.Map{
			"error": "unmarshalling request data",
		})
	}
	admissionReview.Response = &v1beta1.AdmissionResponse{
		Allowed: true,
		UID:     admissionReview.Request.UID,
	}
	containersList := append(pod.Spec.Containers, pod.Spec.InitContainers...)
	message, err := v.imageValidation(containersList)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(&fiber.Map{
			"error": message,
		})
	}
	if message != "" {
		if !v.allowConfig.LogOnly {
			admissionReview.Response.Allowed = false
			admissionReview.Response.Result = &metav1.Status{
				Message: message,
			}
		}
	}
	res, mList := v.resourceValidation(containersList)
	if !res {
		if !v.allowConfig.LogOnly {
			admissionReview.Response.Allowed = false
			admissionReview.Response.Result = &metav1.Status{
				Message: fmt.Sprintf("Resources are undefined for containers: %s", strings.Join(mList[:], ", ")),
			}
		}
	}
	res, err = v.versionValidation(containersList)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(&fiber.Map{
			"error": err.Error(),
		})
	}
	if !res {
		if !v.allowConfig.LogOnly {
			admissionReview.Response.Allowed = false
			admissionReview.Response.Result = &metav1.Status{
				Message: "Version of the image are lower then specified threshold",
			}
		}
	}
	return c.JSON(admissionReview)
}
