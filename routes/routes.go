package routes

import (
	"net/http"
	"project/controllers"
	"project/internal/service"
	"project/middlewares"
	"project/repository"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(db *pgxpool.Pool) *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()

	// Department
	deptRepo := repository.NewDepartmentRepository(db)
	deptService := service.NewDepartmentService(deptRepo)
	deptHandler := controllers.NewDepartmentHandler(deptService)

	api.HandleFunc("/departments", deptHandler.Create).Methods("POST")
	api.HandleFunc("/departments", deptHandler.GetAll).Methods("GET")
	api.HandleFunc("/departments/{id}", deptHandler.GetByID).Methods("GET")
	api.HandleFunc("/departments/{id}", deptHandler.Update).Methods("PUT")
	api.HandleFunc("/departments/{id}", deptHandler.Delete).Methods("DELETE")

	// Employee
	empRepo := repository.NewEmployeeRepository(db)
	empService := service.NewEmployeeService(empRepo)
	empHandler := controllers.NewEmployeeHandler(empService)

	api.HandleFunc("/employees", empHandler.Create).Methods("POST")
	api.HandleFunc("/employees", empHandler.GetAll).Methods("GET")
	api.HandleFunc("/employees/{id}", empHandler.GetByID).Methods("GET")
	api.HandleFunc("/employees/{id}", empHandler.Update).Methods("PUT")
	api.HandleFunc("/employees/{id}", empHandler.Delete).Methods("DELETE")

	// Position
	posRepo := repository.NewPositionRepository(db)
	posSvc := service.NewPositionService(posRepo)
	posHandler := controllers.NewPositionHandler(posSvc)

	api.HandleFunc("/positions", posHandler.Create).Methods("POST")
	api.HandleFunc("/positions", posHandler.GetAll).Methods("GET")
	api.HandleFunc("/positions/{id}", posHandler.GetByID).Methods("GET")
	api.HandleFunc("/positions/{id}", posHandler.Update).Methods("PUT")
	api.HandleFunc("/positions/{id}", posHandler.Delete).Methods("DELETE")

	// Auth
	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthUseCase(userRepo)
	authHandler := controllers.NewAuthHandler(authSvc)

	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	api.Handle("/auth/refresh", middlewares.RefreshMiddleware(http.HandlerFunc(authHandler.Refresh))).Methods("POST")
	api.HandleFunc("/auth/login", authHandler.GetLogin).Methods("GET")

	// Templates
	templateRepo := repository.NewTemplateRepository(db)
	templateService := service.NewTemplateService(templateRepo)
	templateHandler := controllers.NewTemplateHandler(templateService)

	api.HandleFunc("/templates", templateHandler.Create).Methods("POST")
	api.HandleFunc("/templates", templateHandler.GetAll).Methods("GET")
	api.HandleFunc("/templates/{id}", templateHandler.GetByID).Methods("GET")
	api.HandleFunc("/templates/{id}", templateHandler.Update).Methods("PUT")
	api.HandleFunc("/templates/{id}/preview", templateHandler.Preview).Methods("POST")

	sequenceRepo := repository.NewNumberSequenceRepository(db)
	// Documents
	docRepo := repository.NewDocumentRepository(db)
	docService := service.NewDocumentService(docRepo, userRepo, templateRepo, sequenceRepo)
	docHandler := controllers.NewDocumentHandler(docService)

	api.Handle("/documents", middlewares.AuthMiddleware(http.HandlerFunc(docHandler.Create))).Methods("POST")
	api.HandleFunc("/documents", docHandler.GetAll).Methods("GET")
	api.HandleFunc("/documents/{id}", docHandler.GetByID).Methods("GET")
	api.Handle("/getmydocuments", middlewares.AuthMiddleware(http.HandlerFunc(docHandler.GetMyDocs))).Methods("GET")
	api.HandleFunc("/mydocuments", docHandler.GetMyDocumentsPage).Methods("GET")
	api.HandleFunc("/createcertificate", docHandler.GetCreateCertificatePage).Methods("GET")

	//// Usecase-инстансы
	//userUseCase := user.NewUserUseCase(db)
	//productUseCase := product.NewProductUseCase(db)
	//authUseCase := auth.NewAuthUseCase(db)
	//
	//// Handlers
	//userHandler := controllers.NewUserHandler(userUseCase)
	//productHandler := controllers.NewProductHandler(productUseCase)
	//authHandler := controllers.NewAuthHandler(authUseCase)
	//
	//// Auth routes
	//r.Handle("/auth/login", http.HandlerFunc(authHandler.Login)).Methods("POST")
	//r.Handle("/auth/register", http.HandlerFunc(authHandler.Register)).Methods("POST")
	//r.Handle("/auth/refresh", Refresh(http.HandlerFunc(authHandler.Refresh))).Methods("POST")
	//
	//// User routes (только админ может видеть список и менять пользователей)
	//r.Handle("/users", AdminProtected(http.HandlerFunc(userHandler.GetUsers))).Methods("GET")
	//r.Handle("/users/{id}", Protected(http.HandlerFunc(userHandler.GetUserById))).Methods("GET")
	//r.Handle("/users/{id}", AdminProtected(http.HandlerFunc(userHandler.UpdateUser))).Methods("PUT")
	//r.Handle("/users/{id}", AdminProtected(http.HandlerFunc(userHandler.DeleteUser))).Methods("DELETE")
	//
	//// Product routes
	//r.Handle("/products", Protected(http.HandlerFunc(productHandler.GetProducts))).Methods("GET")
	//r.Handle("/products/{id}", Protected(http.HandlerFunc(productHandler.GetProductByID))).Methods("GET")
	//r.Handle("/products", AdminProtected(http.HandlerFunc(productHandler.CreateProduct))).Methods("POST")
	//r.Handle("/products/{id}", AdminProtected(http.HandlerFunc(productHandler.UpdateProduct))).Methods("PUT")
	//r.Handle("/products/{id}", AdminProtected(http.HandlerFunc(productHandler.DeleteProduct))).Methods("DELETE")

	return r
}
