// Copyright 2017 Venkatesh Gopal(vgopal3@jhu.edu), All rights reserved
package main

import (
  "fmt"
  "io/ioutil"
  "math/big"
  crypt "crypto/rand"
  "os"
)

func main() {

  if len(os.Args) != 3 {
    fmt.Println(" \n Follow command line specification \n ./rabin-keygen" +
      "<publickey-file-name> <privatekey-file-name>\n")

  } else {

    publickeyFileName := os.Args[1]
    privateKeyFileName := os.Args[2]

    p := getprimeNumber()
    q := getprimeNumber()

    N := getPublicKey(p,q)

    WritePublicKeyInformationToFile(N,publickeyFileName)
    WritePrivateKeyInformationToFile(N, p,q,privateKeyFileName )

  }

}

func WritePublicKeyInformationToFile(N *big.Int, publickeyFileName string) {

  NStringToWrite := N.String()
  valueToWrite := NStringToWrite

  err := ioutil.WriteFile(publickeyFileName, []byte(valueToWrite), 0644)
  if err != nil {
    fmt.Println("Some Problem in writing to a file")
  }

}

func WritePrivateKeyInformationToFile(N *big.Int, p *big.Int, q *big.Int,
  privateKeyFileName string) {

    NStringToWrite := N.String()
    commaCharacter := ","
    pStringToWrite := p.String()
    qStringToWrite := q.String()

    valueToWrite := NStringToWrite + commaCharacter + pStringToWrite +
    commaCharacter + qStringToWrite

    err := ioutil.WriteFile(privateKeyFileName, []byte(valueToWrite), 0644)
    if err != nil {
      fmt.Println("Some Problem in writing to a file")
    }

}

func getPublicKey(p *big.Int, q *big.Int) (*big.Int) {

  N := big.NewInt(0)
  N = N.Mul(p,q)
  return N
}

func getprimeNumber()(*big.Int) {
    randomNumber := generateNumber()
  // Check for a prime number
  // I'm hardcoding the value of K in primality test to 5
    accuracyFactor := big.NewInt(5);
    resultWhetherPrime := false

    for (!resultWhetherPrime) {
        randomNumber = generateNumber()
        resultWhetherPrime = isaPrimeNumber(randomNumber,accuracyFactor)
        if (resultWhetherPrime) {
          break
        }
      }
      return randomNumber

}


func generateNumber() (*big.Int) {

  n := 64
  b := make([]byte, n)
  _, y := crypt.Read(b)
  if y != nil {
    fmt.Println("Some error")
  }

  z := big.NewInt(0)
  randomNumber := z.SetBytes(b)

  return randomNumber
}

func squareAndMultiple(a *big.Int, b *big.Int, c *big.Int) (*big.Int) {

  // FormatInt will provide the binary representation of a number
  binExp := fmt.Sprintf("%b", b)
  binExpLength := len(binExp)

  initialValue := big.NewInt(0)
  initialValue = initialValue.Mod(a,c)

  // Hold the initial value in result
  result := big.NewInt(0)
  result = result.Set(initialValue)

  // Using the square and multipy algorithm to perform modular exponentation
  for i := 1; i < binExpLength; i++ {

    // 49 is the ASCII representation of 1 and 48 is the ASCII representation
    // of 0
    interMediateResult := big.NewInt(0)
    interMediateResult = interMediateResult.Mul(result,result)
    result = result.Mod(interMediateResult, c)

    if byte(binExp[i]) == byte(49) {
      interResult := big.NewInt(0)
      interResult = interResult.Mul(result,initialValue)
      result = result.Mod(interResult, c)
    }
  }
  return result

}

func isaPrimeNumber(number *big.Int, accuracyFactor *big.Int) (bool) {

  // First finding the value of r, d as per equation ;
  // d * 2pow(r) = n - 1
  if (((big.NewInt(0)).Mod(number,big.NewInt(2))).Cmp(big.NewInt(0)) == 0) {
    // Case where the /dev/urandom has generated an even number
    return false
  } else {

  varNumber := (big.NewInt(0)).Sub(number, big.NewInt(1))

  r := big.NewInt(2)
  // exponentitalR is 2powr(r)
  exponentitalR := big.NewInt(2)

  for true {

    x := big.NewInt(0)
    modValForX := big.NewInt(0)
    x, modValForX = x.DivMod(varNumber, exponentitalR, modValForX)

    if ( (modValForX.Cmp(big.NewInt(0))) == 0) {
    // Fixing value 10000000000 for calculation purpose
    // To resue the squareAndMultiple algorithm but not affect the modulo part
      r = r.Add(r,big.NewInt(1))
      exponentitalR = squareAndMultiple(big.NewInt(2),
      r, big.NewInt(10000000))

      } else {
        break
      }

    }

  r = r.Sub(r,big.NewInt(1))

  exponentitalR = squareAndMultiple(big.NewInt(2),
  r, big.NewInt(10000000))

  d := big.NewInt(0)
  d = d.Div(varNumber,exponentitalR)

  for i := big.NewInt(0); (i.Cmp(accuracyFactor)) == -1;
  i.Add(i,big.NewInt(1)) {

  millerRabinPrimalityTestResult := millerRabinPrimalityTest(number, d,
  r)

  if (millerRabinPrimalityTestResult == false ) {
    return false
      }
    }
    return true
  }
}


func millerRabinPrimalityTest(number *big.Int, d *big.Int,
  r *big.Int) (bool) {

  // As per millerRabinPrimalityTest, we select an "a" in range[2,n-2]
  // Compute a value x = pow(a,d) % number and return true or false
  // based on some checks
  numberTemp := big.NewInt(0)
  numberTemp = (numberTemp.Sub(number, big.NewInt(4)))
  //aTemp := rand.Int63n(numberTemp.Int64()) + 2
  aTemp := int64(100001)
  a := big.NewInt(aTemp)

  x := squareAndMultiple(a,d,number)

  numberMinusOne := (big.NewInt(0)).Sub(number, big.NewInt(1))
  if( ((x.Cmp(big.NewInt(1))) == 0) || ((x.Cmp(numberMinusOne)) == 0)) {
      return true
    }

  loopCount := (big.NewInt(0)).Sub(r,big.NewInt(1))

  for i := big.NewInt(0); (i.Cmp(loopCount)) == -1; i.Add(i,
    big.NewInt(1)) {

    xIntermediate := (big.NewInt(0)).Mul(x,x)

    x = x.Mod(xIntermediate,number)
    if (x.Cmp(big.NewInt(1)) == 0) {
      return false
    }
    if ((x.Cmp(numberMinusOne)) == 0) {
      return true
    }
  }
  return false

}
