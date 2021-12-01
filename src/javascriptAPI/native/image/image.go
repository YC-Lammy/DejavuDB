package image

import (
	"bytes"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strconv"
	"strings"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/blend"
	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/channel"
	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/fcolor"
	"github.com/anthonynsimon/bild/noise"
	"github.com/anthonynsimon/bild/segment"
	"github.com/anthonynsimon/bild/transform"
	"github.com/dop251/goja"
)

func ImageModuleLoader(vm *goja.Runtime, module *goja.Object) {
	e := module.Get("exports")
	export := e.ToObject(vm)

	export.Set("decode", vm.ToValue(func(arg goja.FunctionCall) goja.Value {
		if len(arg.Arguments) == 0 {
			vm.RunString("throw 'TypeError: image.decode: at least 1 argument required, got 0'")
			return goja.Undefined()
		}
		var reader io.Reader
		switch data := arg.Arguments[0].Export().(type) {
		case goja.ArrayBuffer:
			reader = bytes.NewReader(data.Bytes())
		case string:
			reader = bytes.NewReader([]byte(data))
		default:
			if ar, ok := arg.Arguments[0].ToObject(vm).Get("buffer").Export().(goja.ArrayBuffer); ok {
				reader = bytes.NewReader(ar.Bytes())
			} else {
				vm.RunString("throw 'TypeError: image.decode: argument 1 must be ArrayBuffer or String, got " + arg.Arguments[0].ExportType().Name() + "'")
				return goja.Undefined()
			}
		}
		im, format, err := image.Decode(reader)
		if err != nil {
			vm.RunString("throw 'TypeError: image.decode: " + err.Error())
			return goja.Undefined()
		}
		i := &Image{
			vm:           vm,
			image:        im,
			encodeformat: format,
		}
		obj := vm.NewDynamicObject(i)
		i.this = obj
		return obj
	}))

	export.Set("noiseImage", vm.ToValue(func(arg goja.FunctionCall) goja.Value {
		if len(arg.Arguments) < 2 {
			panic(vm.ToValue("TypeError: image.noiseImage: at least 2 arguments required, got " + strconv.Itoa(len(arg.Arguments))))
		}
		i := &Image{
			vm:    vm,
			image: noise.Generate(int(arg.Arguments[0].ToInteger()), int(arg.Arguments[1].ToInteger()), nil),
		}
		obj := vm.NewDynamicObject(i)
		i.this = obj
		return obj
	}))

	export.Set("rect", vm.ToValue(func(c goja.ConstructorCall) *goja.Object {
		if len(c.Arguments) != 4 {
			panic(vm.ToValue("TypeError: image.rect: expected 4 arguments, got " + strconv.Itoa(len(c.Arguments))))
		}
		r := &rect{
			vm:   vm,
			rect: image.Rect(int(c.Arguments[0].ToInteger()), int(c.Arguments[1].ToInteger()), int(c.Arguments[2].ToInteger()), int(c.Arguments[3].ToInteger())),
		}
		obj := vm.NewDynamicObject(r)
		obj.SetPrototype(c.This.Prototype())
		return obj
	}))

	export.Set("image", vm.ToValue(func(c goja.ConstructorCall) *goja.Object {
		var im image.Image
		var imtype string
		var rec image.Rectangle
		var co color.Color
		if len(c.Arguments) > 0 {
			switch arg := c.Arguments[0].Export().(type) {
			case map[string]interface{}:
				var x0, y0, x1, y1 int
				if v, ok := arg["x0"]; ok {
					if v, ok := v.(int64); ok {
						x0 = int(v)
					}
				}
				if v, ok := arg["y0"]; ok {
					if v, ok := v.(int64); ok {
						y0 = int(v)
					}
				}
				if v, ok := arg["x1"]; ok {
					if v, ok := v.(int64); ok {
						x1 = int(v)
					}
				}
				if v, ok := arg["y1"]; ok {
					if v, ok := v.(int64); ok {
						y1 = int(v)
					}
				}
				rec = image.Rect(x0, y0, x1, y1)

				if v, ok := arg["type"]; ok {
					if v, ok := v.(string); ok {
						imtype = v
					}
				}
			case *rect:
				rec = arg.rect
			}
		}
		switch strings.ToLower(imtype) {
		case "alpha":
			im = image.NewAlpha(rec)
		case "alpha16":
			im = image.NewAlpha16(rec)
		case "cmyk":
			im = image.NewCMYK(rec)
		case "gray":
			im = image.NewGray(rec)
		case "gray16":
			im = image.NewGray16(rec)
		case "nrgba":
			im = image.NewNRGBA(rec)
		case "nrgba64":
			im = image.NewNRGBA64(rec)
		case "nycbcra":
			im = image.NewNYCbCrA(rec, image.YCbCrSubsampleRatio410)
		case "rgba":
			im = image.NewRGBA(rec)
		case "rgba64":
			im = image.NewRGBA64(rec)
		case "uniform":
			im = image.NewUniform(co)
		default:
			im = image.NewNRGBA(rec)
		}
		i := &Image{
			vm:    vm,
			image: im,
		}
		obj := vm.NewDynamicObject(i)
		obj.SetPrototype(c.This.Prototype())
		i.this = obj
		return obj
	}))
}

type Image struct {
	this         *goja.Object
	vm           *goja.Runtime
	image        image.Image
	encodeformat string
}

func (im *Image) Get(key string) goja.Value {
	switch key {
	case "clone":
		return im.vm.ToValue(func(arg goja.FunctionCall) goja.Value {
			return im.vm.NewDynamicObject(&Image{
				vm:           im.vm,
				image:        im.image,
				encodeformat: im.encodeformat,
			})
		})
	case "at":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 2 {
				panic(vm.ToValue("TypeError: image.image.at: at least 2 argument required, got " + strconv.Itoa(len(arg.Arguments))))
			}
			co := im.image.At(int(arg.Arguments[0].ToInteger()), int(arg.Arguments[1].ToInteger()))
			if co == nil {
				return goja.Undefined()
			}
			r, g, b, a := co.RGBA()
			obj := vm.NewObject()
			obj.Set("r", r)
			obj.Set("g", g)
			obj.Set("b", b)
			obj.Set("a", a)
			return obj
		})

	case "bound":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			b := im.image.Bounds()
			return vm.NewDynamicObject(&rect{
				vm:   vm,
				rect: b,
			})
		})

	//////////////// transform ////////////////////////////////////////////
	case "crop":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.crop: at least 1 argument required, got 0"))
			}

			var rec image.Rectangle
			if len(arg.Arguments) == 1 {
				switch v := arg.Arguments[0].Export().(type) {
				case *rect:
					rec = v.rect
				default:
					panic(vm.ToValue("TypeError: image.image.rotate: argument 1 must be number or object, got " + arg.Arguments[0].ExportType().Name()))
				}
			} else {
				if len(arg.Arguments) != 4 {
					panic(vm.ToValue("TypeError: image.image.crop: at least 4 argument required, got " + strconv.Itoa(len(arg.Arguments))))
				}
				var x0, y0, x1, y1 int
				x0 = int(arg.Arguments[0].ToInteger())
				y0 = int(arg.Arguments[1].ToInteger())
				x1 = int(arg.Arguments[2].ToInteger())
				y1 = int(arg.Arguments[3].ToInteger())
				rec = image.Rect(x0, y0, x1, y1)
			}

			im.image = transform.Crop(im.image, rec)
			return goja.Undefined()
		})
	case "flipH":
		return im.vm.ToValue(func(arg goja.FunctionCall) goja.Value {
			im.image = transform.FlipH(im.image)
			return goja.Undefined()
		})
	case "flipV":
		return im.vm.ToValue(func(arg goja.FunctionCall) goja.Value {
			im.image = transform.FlipV(im.image)
			return goja.Undefined()
		})
	case "resize":
		return im.vm.ToValue(func(arg goja.FunctionCall) goja.Value {
			var filter = transform.Linear
			if len(arg.Arguments) < 2 {
				panic(im.vm.ToValue("TypeError: image.image.resize: at least 2 argument required, got " + strconv.Itoa(len(arg.Arguments))))
			}
			if len(arg.Arguments) >= 3 {
				switch strings.ToLower(arg.Arguments[2].String()) {
				case "box":
					filter = transform.Box
				case "catmullrom":
					filter = transform.CatmullRom
				case "gaussian":
					filter = transform.Gaussian
				case "lanczos":
					filter = transform.Lanczos
				case "linear":
					filter = transform.Linear
				case "mitchellnetravali":
					filter = transform.MitchellNetravali
				case "nearestneighbor":
					filter = transform.NearestNeighbor
				}
			}

			width := int(arg.Arguments[0].ToInteger())
			height := int(arg.Arguments[1].ToInteger())
			im.image = transform.Resize(im.image, width, height, filter)
			return goja.Undefined()
		})
	case "rotate":
		return im.vm.ToValue(func(arg goja.FunctionCall) goja.Value {

			if len(arg.Arguments) < 1 {
				panic(im.vm.ToValue("TypeError: image.image.rotate: at least 1 argument required, got 0"))
			}

			var rOptions = &transform.RotationOptions{}
			var angle float64

			switch v := arg.Arguments[0].Export().(type) {
			case int64:
				angle = float64(v)
			case float64:
				angle = v
			case string:
				f, err := strconv.ParseFloat(v, 64)
				if err == nil {
					angle = f
				}
			case map[string]interface{}:
				if v, ok := v["angle"]; ok {
					switch v := v.(type) {
					case int64:
						angle = float64(v)
					case float64:
						angle = v
					case string:
						f, err := strconv.ParseFloat(v, 64)
						if err == nil {
							angle = f
						}
					}
				}
				if v, ok := v["pivot"]; ok {
					if v, ok := v.(map[string]interface{}); ok {
						x, ok := v["x"]
						y, ok1 := v["y"]
						if ok && ok1 {
							switch v := x.(type) {
							case float64:
								x = int64(v)
							case int64:
							default:
								ok = false
							}
							if ok {
								switch v := y.(type) {
								case float64:
									y = int64(v)
								case int64:
								default:
									ok = false
								}
								if ok {
									rOptions.Pivot = &image.Point{X: int(x.(int64)), Y: int(y.(int64))}
								}
							}
						}
					}
				}
				if v, ok := v["resizeBounds"]; ok {
					if b, ok := v.(bool); ok {
						rOptions.ResizeBounds = b
					} else {
						switch v := v.(type) {
						case int64:
							if v != 0 {
								rOptions.ResizeBounds = true
							}
						case float64:
							if v != 0 {
								rOptions.ResizeBounds = true
							}
						case string:
							if v != "" {
								rOptions.ResizeBounds = true
							}
						case nil:
							rOptions.ResizeBounds = false
						default:
							rOptions.ResizeBounds = true
						}
					}
				}
			default:
				panic(im.vm.ToValue("TypeError: image.image.rotate: argument 1 must be number or object, got " + arg.Arguments[0].ExportType().Name()))
			}
			if len(arg.Arguments) >= 2 {
				if options, ok := arg.Arguments[1].Export().(map[string]interface{}); ok {
					if v, ok := options["pivot"]; ok {
						if v, ok := v.(map[string]interface{}); ok {
							x, ok := v["x"]
							y, ok1 := v["y"]
							if ok && ok1 {
								switch v := x.(type) {
								case float64:
									x = int64(v)
								case int64:
								default:
									ok = false
								}
								if ok {
									switch v := y.(type) {
									case float64:
										y = int64(v)
									case int64:
									default:
										ok = false
									}
									if ok {
										rOptions.Pivot = &image.Point{X: int(x.(int64)), Y: int(y.(int64))}
									}
								}
							}
						}
					}
					if v, ok := options["resizeBounds"]; ok {
						if b, ok := v.(bool); ok {
							rOptions.ResizeBounds = b
						} else {
							switch v := v.(type) {
							case int64:
								if v != 0 {
									rOptions.ResizeBounds = true
								}
							case float64:
								if v != 0 {
									rOptions.ResizeBounds = true
								}
							case string:
								if v != "" {
									rOptions.ResizeBounds = true
								}
							case nil:
								rOptions.ResizeBounds = false
							default:
								rOptions.ResizeBounds = true
							}

						}

					}
				}
			}

			im.image = transform.Rotate(im.image, angle, rOptions)
			return goja.Undefined()
		})
	case "shearH":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) == 0 {
				panic(vm.ToValue("TypeError: image.image.shearH: at least 1 argument required, got 0"))
			}
			im.image = transform.ShearH(im.image, arg.Arguments[0].ToFloat())
			return goja.Undefined()
		})
	case "shearV":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) == 0 {
				panic(vm.ToValue("TypeError: image.image.shearV: at least 1 argument required, got 0"))
			}
			im.image = transform.ShearV(im.image, arg.Arguments[0].ToFloat())
			return goja.Undefined()
		})
	case "translate":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 2 {
				panic(vm.ToValue("TypeError: image.image.translate: at least 2 argument required, got " + strconv.Itoa(len(arg.Arguments))))
			}
			im.image = transform.Translate(im.image, int(arg.Arguments[0].ToInteger()), int(arg.Arguments[1].ToInteger()))
			return goja.Undefined()
		})

	/////////////// adjust ////////////////////////////////
	case "brightness":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.brightness: at least 1 argument required, got 0"))
			}
			im.image = adjust.Brightness(im.image, arg.Arguments[0].ToFloat())
			return goja.Undefined()
		})
	case "contrast":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.contrast: at least 1 argument required, got 0"))
			}
			im.image = adjust.Contrast(im.image, arg.Arguments[0].ToFloat())
			return goja.Undefined()
		})
	case "gamma":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.gamma: at least 1 argument required, got 0"))
			}
			im.image = adjust.Gamma(im.image, arg.Arguments[0].ToFloat())
			return goja.Undefined()
		})
	case "hue":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.hue: at least 1 argument required, got 0"))
			}
			im.image = adjust.Hue(im.image, int(arg.Arguments[0].ToInteger()))
			return goja.Undefined()
		})
	case "saturation":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.saturation: at least 1 argument required, got 0"))
			}
			im.image = adjust.Saturation(im.image, arg.Arguments[0].ToFloat())
			return goja.Undefined()
		})

	case "blur":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.blur: at least 1 argument required, got 0"))
			}
			radius := arg.Arguments[0].ToFloat()
			if len(arg.Arguments) >= 2 {
				switch strings.ToLower(arg.Arguments[1].String()) {
				case "box":
					im.image = blur.Box(im.image, radius)
				case "gaussian":
					im.image = blur.Gaussian(im.image, radius)
				default:
					im.image = blur.Box(im.image, radius)
				}
				return goja.Undefined()
			}
			im.image = blur.Box(im.image, radius)
			return goja.Undefined()
		})
	case "extractChannels":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.extractChannels: at least 1 argument required, got 0"))
			}
			chanels := []channel.Channel{}
			for _, a := range arg.Arguments {
				switch s := strings.ToLower(a.String()); s {
				case "red", "r":
					chanels = append(chanels, channel.Red)
				case "green", "g":
					chanels = append(chanels, channel.Green)
				case "blue", "b":
					chanels = append(chanels, channel.Blue)
				case "alpha", "a":
					chanels = append(chanels, channel.Alpha)
				default:
					panic(vm.ToValue("TypeError: image.image.extractChannels: expected 'red', 'green', 'blue' or 'alpha', got '" + s + "'"))
				}
			}
			im.image = channel.ExtractMultiple(im.image, chanels...)
			return goja.Undefined()
		})

	case "dilate":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.dilate: at least 1 argument required, got 0"))
			}
			im.image = effect.Dilate(im.image, arg.Arguments[0].ToFloat())
			return goja.Undefined()
		})

	case "edgeDetection":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.edgeDetection: at least 1 argument required, got 0"))
			}
			im.image = effect.EdgeDetection(im.image, arg.Arguments[0].ToFloat())
			return goja.Undefined()
		})
	case "emboss":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			im.image = effect.Emboss(im.image)
			return goja.Undefined()
		})

	case "erode":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.erode: at least 1 argument required, got 0"))
			}
			im.image = effect.Erode(im.image, arg.Arguments[0].ToFloat())
			return goja.Undefined()
		})
	case "grayscale":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.grayscale: at least 1 argument required, got 0"))
			}
			im.image = effect.Grayscale(im.image)
			return goja.Undefined()
		})
	case "invert":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			im.image = effect.Invert(im.image)
			return goja.Undefined()
		})
	case "median":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.median: at least 1 argument required, got 0"))
			}
			im.image = effect.Median(im.image, arg.Arguments[0].ToFloat())
			return goja.Undefined()
		})
	case "sepia":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			im.image = effect.Sepia(im.image)
			return goja.Undefined()
		})
	case "sharpen":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			im.image = effect.Sharpen(im.image)
			return goja.Undefined()
		})
	case "sobel":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			im.image = effect.Sobel(im.image)
			return goja.Undefined()
		})
	case "unsharpMask":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 2 {
				panic(vm.ToValue("TypeError: image.image.unsharpMask: at least 2 argument required, got " + strconv.Itoa(len(arg.Arguments))))
			}
			im.image = effect.UnsharpMask(im.image, arg.Arguments[0].ToFloat(), arg.Arguments[1].ToFloat())
			return goja.Undefined()
		})
	case "segmentThreshold":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.segmentThreshold: at least 1 argument required, got 0"))
			}
			im.image = segment.Threshold(im.image, uint8(arg.Arguments[0].ToInteger()))
			return goja.Undefined()
		})

	case "add":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.add: at least 1 argument required, got 0"))
			}
			if fg, ok := arg.Arguments[0].Export().(*Image); ok {
				im.image = blend.Add(im.image, fg.image)
			} else {
				panic(vm.ToValue("TypeError: image.image.add: argument 1 must be image.image"))
			}
			return goja.Undefined()
		})
	case "blend":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 2 {
				panic(vm.ToValue("TypeError: image.image.blend: at least 2 argument required, got " + strconv.Itoa(len(arg.Arguments))))
			}
			if fg, ok := arg.Arguments[0].Export().(*Image); ok {
				callback, ok := goja.AssertFunction(arg.Arguments[1])
				if ok {
					call := func(c1 fcolor.RGBAF64, c2 fcolor.RGBAF64) fcolor.RGBAF64 {
						a := im.vm.NewObject()
						a.Set("r", c1.R)
						a.Set("g", c1.G)
						a.Set("b", c1.B)
						a.Set("a", c1.A)
						b := im.vm.NewObject()
						b.Set("r", c2.R)
						b.Set("g", c2.G)
						b.Set("b", c2.B)
						b.Set("a", c2.A)
						v, err := callback(im.this, a, b)
						if err != nil {
							vm.RunString("throw 'TypeError: image.image.blend: " + err.Error())
						}
						re := v.ToObject(im.vm)
						return fcolor.RGBAF64{
							R: re.Get("r").ToFloat(),
							G: re.Get("g").ToFloat(),
							B: re.Get("b").ToFloat(),
							A: re.Get("a").ToFloat(),
						}
					}

					im.image = blend.Blend(im.image, fg.image, call)
				} else {
					panic(vm.ToValue("TypeError: image.image.blend: argument 2 must be function"))
				}
			} else {
				panic(vm.ToValue("TypeError: image.image.blend: argument 1 must be image.image"))
			}
			return goja.Undefined()
		})

	case "colorBurn":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.colorBurn: at least 1 argument required, got 0"))
			}
			if fg, ok := arg.Arguments[0].Export().(*Image); ok {
				im.image = blend.ColorBurn(im.image, fg.image)
			} else {
				panic(vm.ToValue("TypeError: image.image.colorBurn: argument 1 must be image.image"))
			}
			return goja.Undefined()
		})
	case "colorDodge":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.colorDodge: at least 1 argument required, got 0"))
			}
			if fg, ok := arg.Arguments[0].Export().(*Image); ok {
				im.image = blend.ColorDodge(im.image, fg.image)
			} else {
				panic(vm.ToValue("TypeError: image.image.colorDodge: argument 1 must be image.image"))
			}
			return goja.Undefined()
		})

	case "darken":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.darken: at least 1 argument required, got 0"))
			}
			if fg, ok := arg.Arguments[0].Export().(*Image); ok {
				im.image = blend.Darken(im.image, fg.image)
			} else {
				panic(vm.ToValue("TypeError: image.image.darken: argument 1 must be image.image"))
			}
			return goja.Undefined()
		})

	case "difference":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.different: at least 1 argument required, got 0"))
			}
			if fg, ok := arg.Arguments[0].Export().(*Image); ok {
				im.image = blend.Difference(im.image, fg.image)
			} else {
				panic(vm.ToValue("TypeError: image.image.difference: argument 1 must be image.image"))
			}
			return goja.Undefined()
		})
	case "divide":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.divide: at least 1 argument required, got 0"))
			}
			if fg, ok := arg.Arguments[0].Export().(*Image); ok {
				im.image = blend.Divide(im.image, fg.image)
			} else {
				panic(vm.ToValue("TypeError: image.image.divide: argument 1 must be image.image"))
			}
			return goja.Undefined()
		})
	case "exclusion":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.exclusion: at least 1 argument required, got 0"))
			}
			if fg, ok := arg.Arguments[0].Export().(*Image); ok {
				im.image = blend.Exclusion(im.image, fg.image)
			} else {
				panic(vm.ToValue("TypeError: image.image.exclusion: argument 1 must be image.image"))
			}
			return goja.Undefined()
		})
	case "lighten":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.lighten: at least 1 argument required, got 0"))
			}
			if fg, ok := arg.Arguments[0].Export().(*Image); ok {
				im.image = blend.Lighten(im.image, fg.image)
			} else {
				panic(vm.ToValue("TypeError: image.image.lighten: argument 1 must be image.image"))
			}
			return goja.Undefined()
		})
	case "linearBurn":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.linearBurn: at least 1 argument required, got 0"))
			}
			if fg, ok := arg.Arguments[0].Export().(*Image); ok {
				im.image = blend.LinearBurn(im.image, fg.image)
			} else {
				panic(vm.ToValue("TypeError: image.image.linearBurn: argument 1 must be image.image"))
			}
			return goja.Undefined()
		})
	case "linearLight":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.linearLight: at least 1 argument required, got 0"))
			}
			if fg, ok := arg.Arguments[0].Export().(*Image); ok {
				im.image = blend.LinearLight(im.image, fg.image)
			} else {
				panic(vm.ToValue("TypeError: image.image.linearLight: argument 1 must be image.image"))
			}
			return goja.Undefined()
		})
	case "multiply":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.multiply: at least 1 argument required, got 0"))
			}
			if fg, ok := arg.Arguments[0].Export().(*Image); ok {
				im.image = blend.Multiply(im.image, fg.image)
			} else {
				panic(vm.ToValue("TypeError: image.image.multiply: argument 1 must be image.image"))
			}
			return goja.Undefined()
		})
	case "normal":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.normal: at least 1 argument required, got 0"))
			}
			if fg, ok := arg.Arguments[0].Export().(*Image); ok {
				im.image = blend.Normal(im.image, fg.image)
			} else {
				panic(vm.ToValue("TypeError: image.image.normal: argument 1 must be image.image"))
			}
			return goja.Undefined()
		})
	case "opacity":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 2 {
				panic(vm.ToValue("TypeError: image.image.opacity: at least 2 argument required, got " + strconv.Itoa(len(arg.Arguments))))
			}
			if fg, ok := arg.Arguments[0].Export().(*Image); ok {
				im.image = blend.Opacity(im.image, fg.image, arg.Arguments[1].ToFloat())
			} else {
				panic(vm.ToValue("TypeError: image.image.opacity: argument 1 must be image.image"))
			}
			return goja.Undefined()
		})
	case "overlay":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.overlay: at least 1 argument required, got 0"))
			}
			if fg, ok := arg.Arguments[0].Export().(*Image); ok {
				im.image = blend.Overlay(im.image, fg.image)
			} else {
				panic(vm.ToValue("TypeError: image.image.overlay: argument 1 must be image.image"))
			}
			return goja.Undefined()
		})
	case "screen":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.screen: at least 1 argument required, got 0"))
			}
			if fg, ok := arg.Arguments[0].Export().(*Image); ok {
				im.image = blend.Screen(im.image, fg.image)
			} else {
				panic(vm.ToValue("TypeError: image.image.screen: argument 1 must be image.image"))
			}
			return goja.Undefined()
		})
	case "softLight":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.softLight: at least 1 argument required, got 0"))
			}
			if fg, ok := arg.Arguments[0].Export().(*Image); ok {
				im.image = blend.SoftLight(im.image, fg.image)
			} else {
				panic(vm.ToValue("TypeError: image.image.softLight: argument 1 must be image.image"))
			}
			return goja.Undefined()
		})
	case "subtract":
		return im.vm.ToValue(func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
			if len(arg.Arguments) < 1 {
				panic(vm.ToValue("TypeError: image.image.subtract: at least 1 argument required, got 0"))
			}
			if fg, ok := arg.Arguments[0].Export().(*Image); ok {
				im.image = blend.Subtract(im.image, fg.image)
			} else {
				panic(vm.ToValue("TypeError: image.image.subtract: argument 1 must be image.image"))
			}
			return goja.Undefined()
		})
	case "sub":
		return im.Get("subtract")
	}
	return goja.Undefined()
}

func (r Image) Set(key string, val goja.Value) bool {
	return false
}

func (r Image) Delete(key string) bool {
	return false
}

func (r Image) Has(key string) bool {
	for _, k := range r.Keys() {
		if k == key {
			return true
		}
	}
	return false
}

func (r Image) Keys() []string {
	return []string{
		"clone",
		"at",
		"bound",
		"crop",
		"flipH",
		"flipV",
		"resize",
		"rotate",
		"shearH",
		"shearV",
		"translate",
		"brightness",
		"contrast",
		"gamma",
		"hue",
		"saturation",
		"blur",
		"extractChannels",
		"dilate",
		"edgeDetection",
		"emboss",
		"erode",
		"grayscale",
		"invert",
		"median",
		"sepia",
		"sharpen",
		"sobel",
		"unsharpMask",
		"segmentThreshold",
		"add",
		"blend",
		"colorBurn",
		"colorDodge",
		"darken",
		"difference",
		"divide",
		"exclusion",
		"lighten",
		"linearBurn",
		"linearLight",
		"multiply",
		"normal",
		"opacity",
		"overlay",
		"screen",
		"softLight",
		"subtract",
		"sub",
	}
}

type rect struct {
	vm   *goja.Runtime
	rect image.Rectangle
	Type string
}

func (r *rect) Get(key string) goja.Value {
	return goja.Undefined()
}

func (r rect) Set(key string, val goja.Value) bool {
	return false
}

func (r rect) Delete(key string) bool {
	return false
}

func (r rect) Has(key string) bool {
	for _, k := range r.Keys() {
		if k == key {
			return true
		}
	}
	return false
}

func (r rect) Keys() []string {
	return []string{}
}
