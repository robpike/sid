// Copyright 2019 Rob Pike. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The sid program reports the current sidereal time. It needs to know the
// terrestrial longitude, which can be provided by a flag or by reading a file in
// the format used by the Plan 9 astro command: one line of text containing the
// latitude, west longitude, and elevation. (Only the longitude is used.)
package main // import "robpike.io/sid"

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"time"
)

var (
	wLong  = flag.Float64("long", 0.0, "west longitude in degrees; default read from sky file")
	sky    = flag.String("sky", "/usr/local/plan9/sky/here", "sky `file` in Plan 9 format")
	julian = flag.Bool("julian", false, "print Julian date as well")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("sid: ")
	flag.Parse()
	// Neat simple algorithm from https://www.aa.quae.nl/en/reken/sterrentijd.html
	t := time.Date(-4713, 11, 24, 12, 0, 0, 0, time.UTC)
	t2000 := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	ΔJ := float64(time.Now().Unix()-t2000.Unix()) / 86400
	const (
		L0 = 99.967794687
		L1 = 360.98564736628603
		L2 = 2.907879E-13
		L3 = -5.302E-22
	)
	var lw = westLongitude()
	θ := L0 + ΔJ*(L1+ΔJ*(L2+L3*ΔJ)) - lw
	θ = math.Mod(θ, 360) / 15
	hours := int(θ)
	θ -= float64(hours)
	θ *= 60
	minutes := int(θ)
	θ -= float64(minutes)
	θ *= 60
	seconds := int(θ)
	fmt.Printf("%.02dh%.02dm%.02ds\n", hours, minutes, seconds)
	if *julian {
		fmt.Printf("Julian date: %.2f\n", float64(time.Now().Unix()-t.Unix())/86400)
	}
}

func westLongitude() float64 {
	if *wLong != 0 {
		return *wLong
	}
	file, err := ioutil.ReadFile(*sky)
	if err != nil {
		log.Fatal("can't read sky file; set -long for longitude")
	}
	var lat, long, elev float64
	if n, _ := fmt.Sscanf(string(file), "%f %f %f", &lat, &long, &elev); n == 3 {
		return long
	}
	log.Fatal("can't parse sky file; set -long for longitude")
	return 0
}
