# ffmcd

ffmcd is a [Golang](https://golang.org) package to generate [ffmpeg](https://ffmpeg.org/) commands by given filterchains.

## Features
* Use another filterchain's output as input programmatically.
* Use input as output directly if there's no filter in the filterchain automatically.

## Limitation
* The generated command is in the format:
  ```bash
  ffmpeg \
  -i "FILE_1"
  -i "FILE_N"
  -filter_complex \
  "
  [0:v]FILTER_1,...FILTER_N[LABEL_for_0_v];
  [0:a]FILTER_1,...FILTER_N[LABEL_for_0_a];
  ......
  [N:v]FILTER_1,...FILTER_N[LABEL_for_N_v];
  [N:a]FILTER_1,...FILTER_N[LABEL_for_N_a];
  [LABEL_for_N_v]FILTER_1,...[OUT_v];
  [LABEL_for_N_a]FILTER_2,...[OUT_a]" \
  -map "[OUT_v]" -map "[OUT_a]" \
  output.mp4
  ```

## Docs
* <https://pkg.go.dev/github.com/northbright/ffmcd>

## License
* [MIT License](LICENSE)
