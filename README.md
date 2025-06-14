# ffmcd

ffmcd is a [Golang](https://golang.org) package to generate [ffmpeg](https://ffmpeg.org/) commands by specifying inputs, filterchains and output.

## Features
* Use another filterchain's output as input programmatically.
* Use input as output directly if there's no filter in the filterchain automatically.

## Limitation
* The generated command is in the following format:
  ```bash
  ffmpeg \
  -i "FILE_1"
  -i "FILE_N"
  -filter_complex \
  "
  [0:v]FILTER_1,FILTER_2,...FILTER_N[0_v];
  [0:a]FILTER_1,FILTER_2,...FILTER_N[0_a];
  [1:v]FILTER_1,FILTER_2,...FILTER_N[1_v];
  [1:a]FILTER_1,FILTER_2,...FILTER_N[1_a];
  ......
  [N:v]FILTER_1,FILTER_2,...FILTER_N[N_v];
  [N:a]FILTER_1,FILTER_2,...FILTER_N[N_a];
  [0_v][0_a][1_v][1_a]...[n_v][n_a]FILTER_HAS_MULTI_INPUT_AND_OUTPUT[out_v][out_a]" \
  -map "[out_v]" -map "[out_a]" \
  output.mp4
  ```

## Docs
* <https://pkg.go.dev/github.com/northbright/ffmcd>

## License
* [MIT License](LICENSE)
