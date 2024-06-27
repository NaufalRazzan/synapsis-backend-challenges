package routes

import (
	"fmt"
	"net/mail"
	"synapsis-backend-test/controllers/product"
	"synapsis-backend-test/controllers/user"
	"synapsis-backend-test/models"
	"synapsis-backend-test/pkg/middlewares"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

func isValidEmail(email string) bool {
	emailAddress, err := mail.ParseAddress(email)

	return err == nil && emailAddress.Address == email
}

func Routes(app *fiber.App) {
	app.Post("/register", register)
	app.Post("/login", login)

	// grouping to apply auth middleware
	productGroup := app.Group("/v1", middlewares.AuthMiddleware())
	productGroup.Get("/listProduct", fetchAllProductsByCategory)
	productGroup.Post("/insertShoppingCart", insertTrx)
	productGroup.Get("/listShoppingCart", fetchShoppingCart)
	productGroup.Delete("/deleteShoppingCart", deleteShoppingCart)
	productGroup.Patch("/checkout", checkoutToPayment)
}

func register(c *fiber.Ctx) error {
	// check method
	if c.Method() != fiber.MethodPost {
		return fiber.NewError(fiber.StatusMethodNotAllowed, "invalid http method")
	}

	// parsing
	var inputData models.Users
	if err := c.BodyParser(&inputData); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "request entity contains invalid or missing data")
	}

	// validation
	validator := validator.New()
	if err := validator.Struct(inputData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "empty fields")
	}

	// check if email is valid
	if !isValidEmail(inputData.Email) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid email address")
	}
	// validate if password is strong password
	if err := passwordvalidator.Validate(inputData.Password, passwordvalidator.GetEntropy("apapapa123#")); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// signup
	res, err := user.SignUp(inputData)
	if err != nil {
		if err.Error() == "user name already exists" {
			return fiber.NewError(fiber.StatusConflict, "user name already exists")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set("Location", fmt.Sprintf("%s/%s", c.BaseURL(), res))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "new user created",
	})
}

func login(c *fiber.Ctx) error {
	// check method
	if c.Method() != fiber.MethodPost {
		return fiber.NewError(fiber.StatusMethodNotAllowed, "invalid http method")
	}

	// parsing
	var inputData models.UsersLogin
	if err := c.BodyParser(&inputData); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "request entity contains invalid or missing data")
	}

	// validation
	validator := validator.New()
	if err := validator.Struct(inputData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "empty fields")
	}

	// signin
	res, err := user.SignIn(inputData.Email, inputData.Password)
	if err != nil {
		if err.Error() == "invalid email or password" {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"acc_token": res,
	})
}

func fetchAllProductsByCategory(c *fiber.Ctx) error {
	// check method
	if c.Method() != fiber.MethodGet {
		return fiber.NewError(fiber.StatusMethodNotAllowed, "invalid http method")
	}

	category_name := c.Query("category")
	if category_name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "query params is required")
	}

	// fetch
	res, err := product.ViewProductsListByCategory(category_name)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if len(res) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": fmt.Sprintf("no products data of %s", category_name),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("fetch %d products based on %s", len(res), category_name),
		"data":    res,
	})
}

func insertTrx(c *fiber.Ctx) error {
	// check method
	if c.Method() != fiber.MethodPost {
		return fiber.NewError(fiber.StatusMethodNotAllowed, "invalid http method")
	}

	// parsing
	var inputData models.InsertTrxpayload
	if err := c.BodyParser(&inputData); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "request entity contains invalid or missing data")
	}

	// validation
	validator := validator.New()
	if err := validator.Struct(inputData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "empty fields")
	}

	res, err := product.InsertProductToShoppingCart(inputData.Product_id, inputData.User_name, inputData.Amount)
	if err != nil {
		if err.Error() == "products does not exists" {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set("location", res)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success insert to shopping cart",
	})
}

func fetchShoppingCart(c *fiber.Ctx) error {
	// check method
	if c.Method() != fiber.MethodGet {
		return fiber.NewError(fiber.StatusMethodNotAllowed, "invalid http method")
	}

	user_name := c.Query("name")
	if user_name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "must provide the name")
	}

	res, err := product.ViewShoppingCartLists(user_name)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if len(res) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "you have not make any transactions",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("fetched %d", len(res)),
		"data":    res,
	})
}

func deleteShoppingCart(c *fiber.Ctx) error {
	// check method
	if c.Method() != fiber.MethodDelete {
		return fiber.NewError(fiber.StatusMethodNotAllowed, "invalid http method")
	}

	trx_id := c.Query("id")
	if trx_id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "must provide an id")
	}

	if err := product.DeleteProductFromShoppingCart(trx_id); err != nil {
		if err.Error() == "no data exists to be deleted" {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "delete success",
	})
}

func checkoutToPayment(c *fiber.Ctx) error {
	// check method
	if c.Method() != fiber.MethodPatch {
		return fiber.NewError(fiber.StatusMethodNotAllowed, "invalid http method")
	}

	trx_id := c.Query("id")
	if trx_id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "must provide an id")
	}

	if err := product.CheckoutToPayment(trx_id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}
