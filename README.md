# Alien-invasion
Aliens are about to invade the earth, this program simulates the invasion and prints the left over cities.

## Prerequisties
If this is your first time encountering Go, please follow the [instructions](https://golang.org/doc/install) to install Go on your computer. This application requires Go 1.17 or above

## Getting Started
Clone the repository and install the dependencies by executing following command
```bash
$ git clone https://github.com/munna0908/alien-invasion.git
```
```bash
cd alien-invasion && make
```

## Usage
### Build
Generate the binary by executing any one of the following commands
```bash
make build
```
or 
```bash
make install 
````
Please ensure that `GOBIN` is added to the system path.

### Run
Start the program by executing following command
```bash
./alieninvasion -aliens <aliens_count> -iterations <max_iterations> -input-file <file_path>
```
In order to stop excution while invasion is in progress,Hit Ctrl-C. This will print the left over cities and gracefully close the program 

Example:
```bash
./alieninvasion -aliens 20 -iterations 10000 -input-file ./file.txt
```
CLI Options
```bash
Usage of ./alieninvasion:
  -aliens int
        Number of aliens
  -input-file string
        Location of input world file
  -iterations int
        Number of iterations (default 10000)
```

### Test
Run the test suite using following command
```bash
make test
```

## Assumptions
- Total No.Of aliens shoulde be <= 2*(No.of cities). 
- No more than two aliens can occupy a city, if a third alien attempts to enter it will be denied.
- Every city should have aleast one neighbour
