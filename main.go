package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"synapsis-backend-test/configs"
	"synapsis-backend-test/pkg/middlewares"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	FRecover "github.com/gofiber/fiber/v2/middleware/recover"
)

func setApp(file *os.File) *fiber.App{
	app := fiber.New()

	// recover middleware
	app.Use(FRecover.New(FRecover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			log.New(io.MultiWriter(os.Stdout, file), "[ERROR]", log.Ldate|log.Ltime).Printf("%s %s: %s\n", c.Path(), c.Method(), e)
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": e,
			})
		},
	}))

	// logger middleware
	app.Use(logger.New(logger.Config{
		Format: "[${time}] [${severity}] ${path} ${method} (${ip}) ${status} ${latency} - ${message}\n",
		CustomTags: map[string]logger.LogFunc{
			"time": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				return output.WriteString(time.Now().Format("2006-01-02 15:04:05"))
			},
			"message": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				if bodybytes := c.Response().Body(); bodybytes != nil{
					var bodyData map[string]interface{}

					err := json.Unmarshal(bodybytes, &bodyData);
					if err == nil{
						msgValue, _ := bodyData["message"].(string)
						return output.WriteString(msgValue)
					} else{
						panic(err)
					}
				}
				return 0, nil
			},
			"severity": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				status := c.Response().StatusCode()

				if status == fiber.StatusInternalServerError{
					return output.WriteString(("ERROR"))
				}
				return output.WriteString("INFO")
			},
		},
		Output: io.MultiWriter(os.Stdout, file),
	}))

	// check valid routes middleware
	app.Use(middlewares.UndefinedRoutesMiddleware())

	// print error middleware for general usecase
	app.Use(middlewares.ErrorMiddleware())

	return app
}

func main(){
	logPath := filepath.Join("logs", "server.log")
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil{
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer file.Close()

	defer func(){
		if r := recover(); r != nil{
			fmt.Fprintln(os.Stderr, r)
			log.New(file, "[ERROR] ", log.Ldate|log.Ltime).Println("Encountered a system error: ", r)
		}
	}()

	app := setApp(file)

	// register routes

	if err := app.Listen(fmt.Sprintf(":%s", configs.GetConfig().Port)); err != nil{
		log.New(file, "[ERROR] ", log.Ldate|log.Ltime).Println("Application failed to start running: ", err)
		os.Exit(1)
	}
}