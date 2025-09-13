package routes

import (
	"net/http"
	"project/controllers"
	"project/internal/service"
	"project/repository"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"

	"project/middlewares"
)

func SetupRoutes(db *pgxpool.Pool) *mux.Router {
	r := mux.NewRouter()

	// Department
	deptRepo := repository.NewDepartmentRepository(db)
	deptService := service.NewDepartmentService(deptRepo)
	deptHandler := controllers.NewDepartmentHandler(deptService)

	r.HandleFunc("/departments", deptHandler.Create).Methods("POST")
	r.HandleFunc("/departments", deptHandler.GetAll).Methods("GET")
	r.HandleFunc("/departments/{id}", deptHandler.GetByID).Methods("GET")
	r.HandleFunc("/departments/{id}", deptHandler.Update).Methods("PUT")
	r.HandleFunc("/departments/{id}", deptHandler.Delete).Methods("DELETE")

	// Employee

	empRepo := repository.NewEmployeeRepository(db)
	empService := service.NewEmployeeService(empRepo)
	empHandler := controllers.NewEmployeeHandler(empService)

	r.HandleFunc("/employees", empHandler.Create).Methods("POST")
	r.HandleFunc("/employees", empHandler.GetAll).Methods("GET")
	r.HandleFunc("/employees/{id}", empHandler.GetByID).Methods("GET")
	r.HandleFunc("/employees/{id}", empHandler.Update).Methods("PUT")
	r.HandleFunc("/employees/{id}", empHandler.Delete).Methods("DELETE")

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

// Middleware обёртки
func Refresh(h http.HandlerFunc) http.Handler {
	return middlewares.RefreshMiddleware(h)
}

func Protected(h http.HandlerFunc) http.Handler {
	return middlewares.AuthMiddleware(h)
}

func AdminProtected(h http.HandlerFunc) http.Handler {
	return middlewares.OnlyAdmin(middlewares.AuthMiddleware(h))
}
