package patch

import (
	"bytes"
	"testing"

	api "github.com/save-abandoned-projects/libgitops/cmd/sample-app/apis/sample"
	"github.com/save-abandoned-projects/libgitops/cmd/sample-app/apis/sample/scheme"
	"github.com/save-abandoned-projects/libgitops/pkg/runtime"
	"github.com/save-abandoned-projects/libgitops/pkg/serializer"
)

var (
	basebytes = []byte(`
{
	"kind": "Car",
	"apiVersion": "sample-app.weave.works/v1alpha1",
	"metadata": {
	  "name": "foo",
	  "uid": "0123456789101112"
	},
	"spec": {
	  "engine": "foo",
	  "brand": "bar"
	}
}`)
	overlaybytes = []byte(`{"status":{"speed":24.7}}`)
)

var carGVK = api.SchemeGroupVersion.WithKind("Car")
var p = NewPatcher(scheme.Serializer)

func TestCreatePatch(t *testing.T) {
	car := &api.Car{
		Spec: api.CarSpec{
			Engine: "foo",
			Brand:  "bar",
		},
	}
	car.SetGroupVersionKind(carGVK)
	b, err := p.Create(car, func(obj runtime.Object) error {
		car2 := obj.(*api.Car)
		car2.Status.Speed = 24.7
		return nil
	})
	if !bytes.Equal(b, overlaybytes) {
		t.Error(string(b), err, car.Status.Speed)
	}
}

func TestApplyPatch(t *testing.T) {
	result, err := p.Apply(basebytes, overlaybytes, carGVK)
	if err != nil {
		t.Fatal(err)
	}
	frameReader := serializer.NewJSONFrameReader(serializer.FromBytes(result))
	if err := scheme.Serializer.Decoder().DecodeInto(frameReader, &api.Car{}); err != nil {
		t.Fatal(err)
	}
}
