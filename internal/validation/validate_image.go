package validation

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ThelonKarrde/Kube-Image-Guardian/internal/config"
	"github.com/ThelonKarrde/Kube-Image-Guardian/internal/rules"
	"k8s.io/api/admission/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Validator struct {
	allowConfig config.AllowConfig
}

func (v *Validator) ReadConfig(path string) {
	v.allowConfig.ReadConfig(path)
}

func (v *Validator) ImageValidation(w http.ResponseWriter, r *http.Request) {
	log.Print(v.allowConfig)
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
	}
	switch r.Method {
	case "POST":
		var admissionReview v1beta1.AdmissionReview
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&admissionReview)
		if err != nil {
			panic(err)
		}
		pod := v1.Pod{}
		if err := json.Unmarshal(admissionReview.Request.Object.Raw, &pod); err != nil {
			http.Error(w, "Error processing image", 500)
		}
		admissionReview.Response = &v1beta1.AdmissionResponse{
			Allowed: true,
			UID:     admissionReview.Request.UID,
		}
		for _, container := range pod.Spec.Containers {
			isLatest, err := rules.IsUsingLatestTag(container.Image)
			if err != nil {
				http.Error(w, "Error processing image", 500)
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
				http.Error(w, "Error processing image", 500)
			}
			if rep {
				break
			}
			reg, err := rules.FromAllowedRegistries(v.allowConfig.AllowedRegistries, container.Image)
			if err != nil {
				http.Error(w, "Error processing image", 500)
			}
			if !reg {
				admissionReview.Response.Allowed = false
				admissionReview.Response.Result = &metav1.Status{
					Message: "Image is from untrusted registry",
				}
				break
			}
		}
		js, err := json.Marshal(admissionReview)
		if err != nil {
			http.Error(w, "Error Marshaling admission review", 500)
		}
		w.Write([]byte(js))
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
