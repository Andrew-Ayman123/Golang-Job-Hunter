package routes

import (
	"net/http"

	"github.com/Andrew-Ayman123/Job-Hunter/internal/handlers"
	"github.com/Andrew-Ayman123/Job-Hunter/internal/middleware"
	"github.com/Andrew-Ayman123/Job-Hunter/internal/services"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	_ "github.com/swaggo/files"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/Andrew-Ayman123/Job-Hunter/docs" // Import the generated docs
)

func SetupRoutes(userHandler *handlers.UserHandler, jwtService *services.JWTService) chi.Router {
	r := chi.NewRouter()

	// Global middleware stack
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(60)) // 60 second timeout
	r.Use(chiMiddleware.Compress(5))

	// API routes
	r.Route("/api/v1", func(v1 chi.Router) {
		setupAPIRoutes(v1, userHandler, jwtService)
	})

	return r
}

func setupAPIRoutes(router chi.Router, userHandler *handlers.UserHandler, jwtService *services.JWTService) {
	// Swagger documentation
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/api/v1/swagger/doc.json"),
	))

	// Health check
	router.Get("/health", healthCheckHandler)

	// General user routes for all users
	setupUserRoutes(router, *userHandler, *jwtService)

}

func setupUserRoutes(router chi.Router, userHandler handlers.UserHandler, jwtService services.JWTService) {
	router.Route("/user", func(user chi.Router) {

		// Public routes (no middleware)
		user.Post("/login", userHandler.HandleUserLogIn) 

		// Protected routes (with middleware)
		user.Group(func(protected chi.Router) {
			protected.Use(middleware.JWTAuth(&jwtService))

			protected.Route("/profile", func(profile chi.Router) {
				profile.Get("/", userHandler.HandleGetProfile) // Get current user's profile
				// User Phone Numbers
				profile.Post("/phone-numbers", userHandler.HandleCreatePhoneNumber)             // Add phone number
				profile.Put("/phone-numbers/{phoneID}", userHandler.HandleUpdatePhoneNumber)    // Update phone number
				profile.Delete("/phone-numbers/{phoneID}", userHandler.HandleDeletePhoneNumber) // Delete phone number

				// User Education
				profile.Post("/education", userHandler.HandleCreateEducation)                 // Add education
				profile.Put("/education/{educationID}", userHandler.HandleUpdateEducation)    // Update education
				profile.Delete("/education/{educationID}", userHandler.HandleDeleteEducation) // Delete education

				// User Experience
				profile.Post("/experience", userHandler.HandleCreateExperience)                  // Add experience
				profile.Put("/experience/{experienceID}", userHandler.HandleUpdateExperience)    // Update experience
				profile.Delete("/experience/{experienceID}", userHandler.HandleDeleteExperience) // Delete experience

				// User Certifications
				profile.Post("/certifications", userHandler.HandleCreateCertification)                     // Add certification
				profile.Put("/certifications/{certificationID}", userHandler.HandleUpdateCertification)    // Update certification
				profile.Delete("/certifications/{certificationID}", userHandler.HandleDeleteCertification) // Delete certification

				// User Projects
				profile.Post("/projects", userHandler.HandleCreateProject)               // Add project
				profile.Put("/projects/{projectID}", userHandler.HandleUpdateProject)    // Update project
				profile.Delete("/projects/{projectID}", userHandler.HandleDeleteProject) // Delete project

				// User Skills
				profile.Post("/skills", userHandler.HandleAddUserSkills)               // Add skill
				profile.Delete("/skills/{skillID}", userHandler.HandleRemoveUserSkill) // Delete skill
			})
		})
	})

	// Public routes (no middleware)
	router.Get("/skills", userHandler.HandleSearchSkills) // Get all skills

	// Applicant specific routes
	setupApplicantRoutes(router, userHandler, jwtService)

	// Recuriter specific routes
	setupRecruiterRoutes(router, userHandler, jwtService)
	// Admin Specific routes
	setupAdminRoutes(router, userHandler, jwtService)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","timestamp":"` +
		r.Context().Value(chiMiddleware.RequestIDKey).(string) + `"}`))
}

func setupApplicantRoutes(router chi.Router, userHandler handlers.UserHandler, jwtService services.JWTService) {
	router.Route("/applicant", func(applicant chi.Router) {
		applicant.Post("/signup", userHandler.HandleApplicantSignUp)

	})
}
func setupRecruiterRoutes(router chi.Router, userHandler handlers.UserHandler, jwtService services.JWTService) {
	router.Route("/recruiter", func(recruiter chi.Router) {
		recruiter.Group(func(protected chi.Router) {
			protected.Use(middleware.JWTAuth(&jwtService))
			protected.Use(middleware.RequireRole("recruiter"))

		})
	})
}
func setupAdminRoutes(router chi.Router, userHandler handlers.UserHandler, jwtService services.JWTService) {
	router.Route("/admin", func(admin chi.Router) {
		admin.Group(func(protected chi.Router) {
			protected.Use(middleware.JWTAuth(&jwtService))
			protected.Use(middleware.RequireRole("admin"))

			// Only admins can create admin and recruiter accounts
			protected.Post("/create-admin", userHandler.HandleCreateAdmin)
			protected.Post("/create-recruiter", userHandler.HandleCreateRecruiter)
			protected.Post("/company", userHandler.HandleCreateCompany)
			protected.Patch("/company/{companyID}", userHandler.HandleUpdateCompany)
			protected.Delete("/company/{companyID}", userHandler.HandleDeleteCompany)
		})
	})
}
