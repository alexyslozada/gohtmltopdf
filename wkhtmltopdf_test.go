package gohtmltopdf

import (
	"bytes"
	"context"
	"testing"
)

const html = `
<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Document</title>
</head>
<body>
    <h1 style="font-family: sans-serif; font-size: 4rem;">Hola mundo</h1>
    <p style="font-family: serif; font-size: 0.8rem;">Lorem ipsum dolor sit amet, consectetur adipisicing elit. Consequuntur doloribus nam nemo odio quam quo repudiandae sunt tenetur. Ab aliquam aut beatae dicta minus nesciunt quam quidem similique, temporibus vero.</p>
</body>
</html>
`

func Test_run(t *testing.T) {
	src := bytes.NewBufferString(html)
	gen := NewGenerator(src)
	data, err := gen.run(context.Background())
	if err != nil {
		t.Fatalf("Got an unexpected error generating pdf: %v", err)
	}

	err = writeFile(data)
	if err != nil {
		t.Fatalf("Got an unexepected error writing file: %v", err)
	}
}
