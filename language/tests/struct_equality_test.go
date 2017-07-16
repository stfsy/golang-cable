package tests

import (
	"testing"

	assert "github.com/stfsy/golang-assert"
)

type Person struct {
	name string
	age  int
}

type Vehicle struct {
	wheels  int
	mileage int64
}

var paul = Person{
	name: "Paul",
	age:  32}

var pepe = Person{
	name: "Pepe",
	age:  64}

var car = Vehicle{
	wheels: 4}

var sedan = Vehicle{
	wheels: 4}

var bike = Vehicle{
	wheels: 2}

func TestPaulEqualsPaul(t *testing.T) {
	assert.Equal(t, paul, paul, "Paul is paul")
}

func TestPaulDoesNotEqualPepe(t *testing.T) {
	assert.NotEqual(t, paul, pepe, "Paul is not pepe")
}

func TestCarEqualsSedan(t *testing.T) {
	assert.Equal(t, car, sedan, "Car is a Sedan")
}

func TestCarDoesNotEqualBike(t *testing.T) {
	assert.NotEqual(t, car, bike, "Car is not a bike")
}
