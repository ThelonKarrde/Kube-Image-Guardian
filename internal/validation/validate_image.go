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
			http.Error(w, "Error decoding request data", 500)
			v.errorLog.Println("Error decoding request data")
		}
		pod := v1.Pod{}
		err = json.Unmarshal(admissionReview.Request.Object.Raw, &pod)
		if err != nil {
			http.Error(w, "Error unmarshalling request data", 500)
			v.errorLog.Println("Error unmarshalling request data")
		}
		admissionReview.Response = &v1beta1.AdmissionResponse{
			Allowed: true,
			UID:     admissionReview.Request.UID,
		}
		for _, container := range pod.Spec.Containers {
			isLatest, err := rules.IsUsingLatestTag(container.Image)
			if err != nil {
				http.Error(w, "Error processing image name", 500)
				v.errorLog.Printf("Error normalizing image name: %s", container.Image)
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
				v.errorLog.Printf("Error normalizing image name: %s", container.Image)
			}
			if !rep {
				reg, err := rules.FromAllowedRegistries(v.allowConfig.AllowedRegistries, container.Image)
				if err != nil {
					http.Error(w, "Error processing image", 500)
					v.errorLog.Printf("Error normalizing image name: %s", container.Image)
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
		js, err := json.Marshal(admissionReview)
		if err != nil {
			http.Error(w, "Error Marshaling admission review response data", 500)
			v.errorLog.Print("Error Marshaling admission review response data")
		}
		w.Write([]byte(js))
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
