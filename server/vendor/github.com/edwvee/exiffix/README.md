Exiffix
======
Exiffix is a one function golang library made to be a replacement for image.Decode to handle orientation stored in EXIF data.

```exiffix.Decode``` has mostly the same signature as ```image.Decode```. The difference is that it requires ```io.ReadSeeker``` instead of ```io.Seeker```.

Example
======
```go
file, err := os.Open(path)
if err != nil{
  panic(err)
}
defer file.Close()
img, fmt, err := exiffix.Decode(file)
```

Special thanks to Macilias: part of code is taken from his [repo](https://github.com/Macilias/go-images-orientation)
