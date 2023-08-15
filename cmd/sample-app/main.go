package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/save-abandoned-projects/libgitops/cmd/common"
	"github.com/save-abandoned-projects/libgitops/cmd/sample-app/apis/sample/scheme"
	"github.com/save-abandoned-projects/libgitops/cmd/sample-app/apis/sample/v1alpha1"
	"github.com/save-abandoned-projects/libgitops/pkg/logs"
	"github.com/save-abandoned-projects/libgitops/pkg/runtime"
	"github.com/save-abandoned-projects/libgitops/pkg/serializer"
	"github.com/save-abandoned-projects/libgitops/pkg/storage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var manifestDirFlag = pflag.String("data-dir", "/tmp/libgitops/manifest", "Where to store the YAML files")

func main() {
	// Parse the version flag
	common.ParseVersionFlag()

	// Run the application
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	// Create the manifest directory
	if err := os.MkdirAll(*manifestDirFlag, 0755); err != nil {
		return err
	}

	// Set the log level
	logs.Logger.SetLevel(logrus.InfoLevel)

	plainStorage := storage.NewGenericStorage(
		storage.NewGenericRawStorage(*manifestDirFlag, v1alpha1.SchemeGroupVersion, serializer.ContentTypeYAML),
		scheme.Serializer,
		[]runtime.IdentifierFactory{runtime.Metav1NameIdentifier},
	)
	defer func() { _ = plainStorage.Close() }()

	e := common.NewEcho()

	e.GET("/plain/:name", func(c echo.Context) error {
		name := c.Param("name")
		if len(name) == 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "Please set name")
		}

		obj, err := plainStorage.Get(common.CarKeyForName(name))
		if err != nil {
			return err
		}
		var content bytes.Buffer
		if err := scheme.Serializer.Encoder().Encode(serializer.NewJSONFrameWriter(&content), obj); err != nil {
			return err
		}
		return c.JSONBlob(http.StatusOK, content.Bytes())
	})

	e.POST("/plain/:name", func(c echo.Context) error {
		name := c.Param("name")
		if len(name) == 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "Please set name")
		}

		if err := plainStorage.Create(common.NewCar(name)); err != nil {
			return err
		}
		return c.String(200, "OK!")
	})

	e.PUT("/plain/:name", func(c echo.Context) error {
		name := c.Param("name")
		if len(name) == 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "Please set name")
		}

		if err := common.SetNewCarStatus(plainStorage, common.CarKeyForName(name)); err != nil {
			return err
		}
		return c.String(200, "OK!")
	})

	return common.StartEcho(e)
}
