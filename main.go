package main

import (
  "log"
  "image"
  "image/draw"
  "time"
  "os/exec"
  "periph.io/x/periph/host"
  "periph.io/x/periph/conn/i2c/i2creg"
  "periph.io/x/periph/devices/ssd1306/image1bit"
  "periph.io/x/periph/devices/ssd1306"
  "golang.org/x/image/font"
  "golang.org/x/image/font/basicfont"
  "golang.org/x/image/math/fixed"
)

func main() {

  // Make sure periph is initialized.
  if _, err := host.Init(); err != nil {
      log.Fatal(err)
  }

  // Use i2creg I²C bus registry to find the first available I²C bus.
  b, err := i2creg.Open("")
  if err != nil {
      log.Fatal(err)
  }
  defer b.Close()

  var opts = ssd1306.Opts{
    W:             128,
    H:             64,
    Rotated:       true,
    Sequential:    false,
    SwapTopBottom: false,
  }

  dev, err := ssd1306.NewI2C(b, &opts)
  if err != nil {
      log.Fatalf("failed to initialize ssd1306: %v", err)
  }

  // Draw on it.
  img := image1bit.NewVerticalLSB(dev.Bounds())


  f := basicfont.Face7x13
  // f.Height = 10
  // f.Ascent = 8
  // f.Width = 4
  // f.Advance = 5
  // f.Mask =
  line1 := font.Drawer{
    Dst:  img,
    Src:  &image.Uniform{image1bit.On},
    Face: f,
    Dot:  fixed.P(0, 13),
  }
  line2 := font.Drawer{
    Dst:  img,
    Src:  &image.Uniform{image1bit.On},
    Face: f,
    Dot:  fixed.P(0, 13*2),
  }
  line3 := font.Drawer{
    Dst:  img,
    Src:  &image.Uniform{image1bit.On},
    Face: f,
    Dot:  fixed.P(0, 13*3),
  }
  line4 := font.Drawer{
    Dst:  img,
    Src:  &image.Uniform{image1bit.On},
    Face: f,
    Dot:  fixed.P(0, 13*4),
  }

  for {
    draw.Draw(img, dev.Bounds(), &image.Uniform{image1bit.Off}, image.ZP, draw.Src)
    img.Set(opts.W, opts.H, &image.Uniform{image1bit.On})
    line1.Dot = fixed.P(0, 13)
    line2.Dot = fixed.P(0, 13*2)
    line3.Dot = fixed.P(0, 13*3)
    line4.Dot = fixed.P(0, 13*4)

    out, _ := exec.Command("bash", "-c", "top -bn1 | grep load | awk '{printf \"%.2f\", $(NF-2)}'").Output()
    line1.DrawString("LA: ")
    line1.DrawBytes(out)
    out, _ = exec.Command("bash", "-c", "cat /sys/class/thermal/thermal_zone0/temp | awk '{printf \"%.0f\", $1/1000}'").Output()
    line1.DrawString(" TMP: ")
    line1.DrawBytes(out)

    out, _ = exec.Command("bash", "-c", "hostname -I | cut -d ' ' -f1 | tr -d '\n'").Output()
    line2.DrawString("IP: ")
    line2.DrawBytes(out)

    out, _ = exec.Command("bash", "-c", "free -m | awk 'NR==2{printf \"%s/%sMB %.0f%%\", $3,$2,$3*100/$2 }'").Output()
    line3.DrawString("MEM: ")
    line3.DrawBytes(out)
    out, _ = exec.Command("bash", "-c", "df -h | awk '$NF==\"/\"{printf \"%d/%dGB %s\", $3,$2,$5}'").Output()
    line4.DrawString("SD: ")
    line4.DrawBytes(out)

    if err = dev.Draw(dev.Bounds(), img, image.Point{}); err != nil {
      log.Fatal(err)
    }
    time.Sleep(5000000 * time.Nanosecond)
  }
}
