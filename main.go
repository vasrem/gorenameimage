package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type byDateCreated []os.FileInfo

// Sorting functions
func (file byDateCreated) Len() int {
	return len(file)
}

func (file byDateCreated) Swap(i, j int) {
	file[i], file[j] = file[j], file[i]
}

func (file byDateCreated) Less(i, j int) bool {
	date1 := getDate(file[i])
	date2 := getDate(file[j])
	return date1.Before(date2)
}

// Get date of image creation
func getDate(file os.FileInfo) time.Time {
	f, err := os.Open(file.Name())
	if err != nil {
		log.Fatal(err)
	}
	x, err := exif.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	date, _ := x.DateTime()
	return date
}

// Copy file from src to dst
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

// Check if path exists
func exists(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	if f.Mode().IsRegular() {
		return false, nil
	}
	return true, nil
}

// Ask for confirmation function
func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

// Closure to generate image name
func imageNameGenerator(prefix string) func() string {
	i := 0
	return func() string {
		i++
		return prefix + "_" + fmt.Sprintf("%06v", i)
	}
}

func main() {

	input := flag.String("input", "", "Directory that contains the images. Provide absolute path.")
	output := flag.String("output", "", "Directory that will contain the renamed images. Provide absolute path.")
	mode := flag.String("mode", "copy", "copy or move")
	prefix := flag.String("prefix", "DSC", "Desired prefix for names")
	flag.Parse()

	// Check mode
	if *mode != "copy" && *mode != "move" {
		fmt.Println("Wrong mode.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Check input
	if *input == "" {
		fmt.Println("You have to insert the absolute path of input directory")
		flag.PrintDefaults()
		os.Exit(1)
	}
	ex, _ := exists(*input)
	if ex != true {
		fmt.Println("Wrong input path.")
		os.Exit(1)
	}

	// Check output
	if *output == "" {
		fmt.Println("You have to insert the absolute path of output directory")
		flag.PrintDefaults()
		os.Exit(1)
	}
	ex, _ = exists(*output)
	if ex != true {
		fmt.Printf("Wrong output path.")
		os.Exit(1)
	}

	// Check if sb wants to copy into same folder
	if *output == *input && *mode == "copy" {
		fmt.Printf("Please use flag -mode=move to rename images at input directory")
		os.Exit(1)
	}

	// Ask confirmation of sb's command
	c := askForConfirmation("Are you sure that you want to " + *mode + " all files from " + *input + " to " + *output + " ?")
	if c == false {
		os.Exit(1)
	}

	// find files in *input path
	ext := "(?i).jpg$"
	var files []os.FileInfo
	filepath.Walk(*input, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(ext, f.Name())
			if err == nil && r {
				files = append(files, f)
			}
		}
		return nil
	})

	os.Chdir(*input)

	// sort by date created
	sort.Sort(byDateCreated(files))

	os.Chdir(*output)

	// do the job
	name := imageNameGenerator(*prefix)
	for _, file := range files {
		in := *input + "/" + file.Name()
		suffix := filepath.Ext(file.Name())
		out := *output + "/" + name() + suffix
		if *mode == "move" {
			os.Rename(in, out)
		} else {
			copyFile(in, out)
		}
	}
}
