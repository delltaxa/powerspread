package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"
)

func randkee() int {
	rand.Seed(time.Now().UTC().UnixNano())

	return rand.Intn(2222) + 100
}

func ghex(s string) string {
	// Create a new bytes buffer
	buf := new(bytes.Buffer)

	// Write the string to the buffer
	buf.WriteString(s)

	// Return the hexadecimal representation of the buffer
	return hex.EncodeToString(buf.Bytes())
}

func rufus(s string, rkee int) string {
	var result string

	for _,i := range s {
			result += string(rune(int(i)+rkee))
	}

	return result
}

func benc(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func pogofuscate(payload string) string {
	var rkee int = randkee()

	var result string = `$______="SEVKsHgHWhRiGzRNahXN";Set-Alias -Name ___ -Value Set-Alias;___ -Name ____ -Value ("{1}{2}{0}"-f $______[18],$______[11],$______[1]);___ -Name __ -Value ____;$r=""; [System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String("`+benc(rufus(ghex(payload), rkee))+`")).ToCharArray()|ForEach-Object{$x=[int]$_;$n=$x-`+strconv.Itoa(rkee)+`;$nx=[char]$n;$r=$r+$nx};-join("$r"-split'(..)'|?{$_}|%%{[char][convert]::ToUInt32($_,16)})|&("_"+"_")`
	return result
}
