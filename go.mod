module training

go 1.13

require (
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/google/uuid v1.1.2
	github.com/gorilla/mux v1.8.0
	github.com/sirupsen/logrus v1.7.0
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83 // indirect
	gorm.io/driver/postgres v1.0.8
	gorm.io/gorm v1.20.12
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v1.0.0 // indirect
	sigs.k8s.io/controller-runtime v0.8.3
)
