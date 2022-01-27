package validation

import (
	"encoding/json"
	"log"

	"github.com/ThelonKarrde/Kube-Image-Guardian/internal/config"
	"github.com/ThelonKarrde/Kube-Image-Guardian/internal/rules"
	"github.com/gofiber/fiber/v2"
	"k8s.io/api/admission/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (v *Validator) ImageValidation(c *fiber.Ctx) error {
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
	for _, container := range pod.Spec.Containers {
		isLatest, err := rules.IsUsingLatestTag(container.Image)
		if err != nil {
			v.errorLog.Printf("Error normalizing image name: %s", container.Image)
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(&fiber.Map{
				"error": "processing image name",
			})
		}
		if !v.allowConfig.AllowLatest && isLatest {
			admissionReview.Response.Allowed = false
			admissionReview.Response.Result = &metav1.Status{
				Message: "Latest tag is not allowed",
			}
			break
		}
		rep, err := rules.FromAllowedRepository(v.allowConfig.AllowedRepositories, container.Image)
		if err != nil {
			v.errorLog.Printf("Error normalizing image name: %s", container.Image)
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(&fiber.Map{
				"error": "processing image name",
			})
		}
		if !rep {
			reg, err := rules.FromAllowedRegistries(v.allowConfig.AllowedRegistries, container.Image)
			if err != nil {
				v.errorLog.Printf("Error normalizing image name: %s", container.Image)
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(&fiber.Map{
					"error": "processing image name",
				})
			}
			if !reg {
				admissionReview.Response.Allowed = false
				admissionReview.Response.Result = &metav1.Status{
					Message: "Image is from untrusted registry",
				}
				break
			}
		}
	}
	return c.JSON(admissionReview)
}
