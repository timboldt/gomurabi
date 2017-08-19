package kingdomstate

import (
    . "github.com/go-check/check"
    "testing"
	"fmt"
    "math/rand"
	"time"
)

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type S struct{}
var _ = Suite(&S{})
var randgen *rand.Rand

func (s *S) SetUpTest(c *C) {
	randgen = rand.New(rand.NewSource(time.Now().UnixNano()))	
}

// -----------------------------------------------------------------------
// IntegerBetween checker.

type integerBetweenChecker struct {
	*CheckerInfo
}

// The IntegerBetween checker verifies that the obtained value is >= to
// the first expected value, and <= the second expected value
// according to usual Go uint semantics for >=, <=.
//
// For example:
//
//     c.Assert(value, IntegerBetween, 42, 45)
//
var IntegerBetween Checker = &integerBetweenChecker{
	&CheckerInfo{Name: "IntegerBetween", Params: []string{"obtained", "expected >=", "expected <="}},
}

func (checker *integerBetweenChecker) Check(params []interface{}, names []string) (result bool, error string) {
	defer func() {
		if v := recover(); v != nil {
			result = false
			error = fmt.Sprint(v)
		}
	}()
	return params[0].(int) >= params[1].(int) && params[0].(int) <= params[2].(int), ""
}

func (s *S) TestRandomPricePerAcre(c *C) {
	c.Check(int(RandomPricePerAcre(nil)), Equals, 21)
	
	totals := make(map[uint]int, 100)
	for i:=0; i<1000; i++ {
		res := RandomPricePerAcre(randgen)
		totals[res] = totals[res] + 1
	}
	c.Check(totals[16] == 0, Equals, true)
	for i:=uint(17); i<27; i++ {
		c.Check(totals[i] != 0, Equals, true)		
	}
	c.Check(totals[27] == 0, Equals, true)
}

func (s *S) TestRandomYieldPerAcre(c *C) {
	c.Check(int(RandomYieldPerAcre(nil)), Equals, 3)
	
	totals := make(map[uint]int, 100)
	for i:=0; i<1000; i++ {
		res := RandomYieldPerAcre(randgen)
		totals[res] = totals[res] + 1
	}
	c.Check(totals[0] == 0, Equals, true)
	for i:=uint(1); i<6; i++ {
		c.Check(totals[i] != 0, Equals, true)		
	}
	c.Check(totals[6] == 0, Equals, true)
}


func (s *S) TestRandomRatPercent(c *C) {
	c.Check(int(RandomRatPercent(nil)), Equals, 10)
	
	totals := make(map[uint]int, 100)
	for i:=0; i<1000; i++ {
		res := RandomRatPercent(randgen)
		totals[res] = totals[res] + 1
	}

	c.Check(totals[0] != 0, Equals, true)
	c.Check(totals[9] == 0, Equals, true)
	for i:=uint(10); i<31; i++ {
		c.Check(totals[i] != 0, Equals, true)		
	}
	c.Check(totals[31] == 0, Equals, true)
}

func (s *S) TestRandomPlagueHappened(c *C) {
	c.Check(RandomPlagueHappened(nil, 1), Equals, false)
	c.Check(RandomPlagueHappened(nil, 2), Equals, false)
	c.Check(RandomPlagueHappened(nil, 3), Equals, false)
	c.Check(RandomPlagueHappened(nil, 4), Equals, true)
	c.Check(RandomPlagueHappened(nil, 9), Equals, false)

	seenTrue := false
	seenFalse := false
	for i:=uint(0); i<1000; i++ {
		res := RandomPlagueHappened(randgen, i%10 + 1)
		if (res) { seenTrue = true } else { seenFalse = true }
	}
	c.Check(seenTrue, Equals, true)
	c.Check(seenFalse, Equals, true)
}

func (s *S) Test_min(c *C) {
	c.Check(min(3,4) == 3, Equals, true)
	c.Check(min(99,0) == 0, Equals, true)
	c.Check(min(43, 43) == 43, Equals, true)
}

func (s *S) Test_max(c *C) {
	c.Check(max(3,4) == 4, Equals, true)
	c.Check(max(99,0) == 99, Equals, true)
	c.Check(max(43, 43) == 43, Equals, true)
}

// Deterministic sequence helper functions
func acresToBuy(year uint, acreage uint) uint {
	if year%3 == 1 { return year * acreage / 100 }
	return 0
}
func acresToSell(year uint, acreage uint) uint {
	if year%3 == 2 { return year * acreage / 100 }
	return 0
}
func grainForFood(year uint, population uint) uint { return (100 - year * 2) * population * GrainPerPerson / 100 }
func acresToPlant(year uint) uint { return 10000 }

func (s *S) TestDeterministicSequence(c *C) {
	expectedEOYPopulation := []uint{100,98,94,88,49,44,38,32,33,27,21}
	expectedEOYAcreage := []uint{1000,1010,990,990,1029,978,978,1046,963,963,1059}
	expectedEOYGrain := []uint{2800,2840,3470,3767,3527,5547,6289,5509,7499,7749,5997}
	expectedEOYStillInOffice := []bool{true, true, true, true, true, true, true, true, true, true, false}
	
	var ks KingdomState
	var year uint
	
	ks.SetupInitialState(nil)
	for year = 1; year<=10; year++ {
		ks.TallyUpYear(acresToBuy(year, ks.acreage), acresToSell(year, ks.acreage), grainForFood(year, ks.population), acresToPlant(year))
//fmt.Printf("%v\n", ks)      
		c.Check(ks.yearOfRule, Equals, year)
		c.Check(int(ks.population), Equals, int(expectedEOYPopulation[year]))
//fmt.Printf("%d %d\n", expectedEOYAcreage[year], expectedEOYGrain[year])
		c.Check(int(ks.acreage), Equals, int(expectedEOYAcreage[year]))
		c.Check(int(ks.grain), Equals, int(expectedEOYGrain[year]))
		c.Check(ks.stillInOffice, Equals, expectedEOYStillInOffice[year])
	}
}

func Benchmark_1(b * testing.B) {
	randgen = rand.New(rand.NewSource(time.Now().UnixNano()))
	for i:=0; i<100000; i++ {
		var ks KingdomState

		ks.SetupInitialState(randgen)
		for ks.StillInOffice() {
			ks.TallyUpYear(0, 50, 2000, 10)
		}
	}
}
/*
//////////////


type KingdomState struct {
	deterministic bool
	stillInOffice bool
	
	yearOfRule uint
	population uint
	acreage uint
	grain uint
	
	pricePerAcre uint
	
	harvestPerAcre uint
	percentEatenByRats uint
	plagueHappened bool
	nextYearPricePerAcre uint
}

func (ks *KingdomState) SetupInitialState(deterministic bool) {
	ks.deterministic = deterministic
	ks.stillInOffice = true
	
	ks.population = 100
	ks.acreage = 1000
	ks.grain = 2800
	
	ks.percentEatenByRats = 10
	ks.nextYearPricePerAcre = RandomPricePerAcre(ks.deterministic)
}

func (ks *KingdomState) TallyUpYear(acresToBuy, acresToSell, grainForFood, acresToPlant uint) {
	if !ks.stillInOffice { panic(0) }
	
	var startOfYearPopulation uint = ks.population
	var starvationVictims uint
	
	// Cummulative values
	ks.yearOfRule++
	ks.pricePerAcre = ks.nextYearPricePerAcre

	// Random events
	ks.harvestPerAcre = RandomYieldPerAcre(ks.deterministic)
	ks.percentEatenByRats = RandomRatPercent(ks.deterministic)
	ks.plagueHappened = RandomPlagueHappened(ks.deterministic, ks.yearOfRule)  
	ks.nextYearPricePerAcre = RandomPricePerAcre(ks.deterministic)
	
	// Buy land
	grainUsedToBuyLand := min(acresToBuy * ks.pricePerAcre, ks.grain)
	ks.grain -= grainUsedToBuyLand
	ks.acreage += grainUsedToBuyLand / ks.pricePerAcre
	
	// Sell land
	grainFromSaleOfLand := min(acresToSell, ks.acreage) * ks.pricePerAcre
	ks.grain += grainFromSaleOfLand
	ks.acreage -= grainFromSaleOfLand / ks.pricePerAcre
	
	// Feed the people
	peopleFed := min(ks.grain / GrainPerPerson, ks.population)
	ks.grain -= peopleFed * GrainPerPerson
	
	// Plant the fields, harvest them, and deal with the rats
	acresForPlanting :=	min( min(acresToPlant, ks.population * AcresPerPerson), ks.acreage )
	grainPlanted := min(ks.grain, acresForPlanting / AcresPerBushel)
	acresPlanted := grainPlanted * AcresPerBushel
	ks.grain += acresPlanted * ks.harvestPerAcre
	ks.grain -= ks.percentEatenByRats * ks.grain / 100
	
	// Adjust population counts
	if ks.plagueHappened { ks.population = ks.population / 2}

	if ks.population > peopleFed {
		// Starvation occurs if not everyone was fed
		starvationVictims = ks.population - peopleFed
		ks.population -= starvationVictims
	} else if ks.population > 0 {
		// Allow immigrants if nobody starved and there are still people around
		ks.population += (20 * ks.acreage + ks.grain) / (100 * ks.population) + 1
	}

	// Determine if the game is over
	ks.stillInOffice = (
		ks.yearOfRule < 10 &&
		ks.population > 0 &&
		starvationVictims < 45 * startOfYearPopulation)
}

class KingdomStateTest extends FunSuite with ShouldMatchers {

  val randomInitialState = new KingdomState()
  val deterministicInitialState = new KingdomState(deterministic = true)

  var deterministicStates = new Array[KingdomState](11)

  def acresToBuyOrSell(i: Int): Int = {
    (i % 3) match {
      case 0 => 0
      case 1 => i * deterministicStates(i - 1).endOfYearAcreage / 100 
      case 2 => -1 * i * deterministicStates(i - 1).endOfYearAcreage / 100
    }
  }

  def grainForFood(i: Int): Int = {
    (100 - i * 2) * deterministicStates(i - 1).endOfYearPopulation * KingdomState.grainPerPerson / 100 
  }

  def acresToPlant(i: Int): Int = {
    10000
    //deterministicStates(i - 1).endOfYearAcreage + acresToBuyOrSell(i)
  }
  
  for (i <- 0 to 10) {
    deterministicStates(i) =
      if (i == 0) deterministicInitialState
      else new KingdomState(deterministicStates(i - 1), acresToBuyOrSell = acresToBuyOrSell(i), grainForFood = grainForFood(i), acresToPlant = acresToPlant(i)) 
      
    val d = deterministicStates(i)
    println("Year: %d Pop: %d (S=%d P=%d I=%d) Grain: %d (B=%d S=%d F=%d P=%d)".format(
      d.yearOfRule,
      d.endOfYearPopulation, d.starvationVictims, d.plagueVictims, d.immigrants,
      d.endOfYearGrain, d.grainUsedToBuyLand, d.grainFromSaleOfLand, d.grainFedToPeople, d.grainPlanted
      )
    )
  }
  
  val expectedEOYPopulation = Array(100,98,94,88,49,44,38,32,33,27,21)
  val expectedEOYAcreage = Array(1000,1010,990,990,1029,978,978,1046,963,963,1059)
  val expectedEOYGrain = Array(2800,2840,3470,3767,3527,5547,6289,5509,7499,7749,5997)
  val expectedEOYStillInOffice = Array(true, true, true, true, true, true, true, true, true, true, false)
  
  test("yearOfRule") {
    randomInitialState.yearOfRule should be (0)
    for (i <- 0 to 10)
      deterministicStates(i).yearOfRule should be (i)
  }

  test("startOfYearPopulation") {
    randomInitialState.startOfYearPopulation should be >= (0)
    for (i <- 0 to 10) {
      deterministicStates(i).startOfYearPopulation should be >= (0)
      if (i > 0) deterministicStates(i).startOfYearPopulation should be (deterministicStates(i - 1).endOfYearPopulation)
    }
  }

  test("startOfYearAcreage") {
    randomInitialState.startOfYearAcreage should be >= (0)
    for (i <- 0 to 10) {
      deterministicStates(i).startOfYearAcreage should be >= (0)
      if (i > 0) deterministicStates(i).startOfYearAcreage should be (deterministicStates(i - 1).endOfYearAcreage)
    }
  }

  test("startOfYearGrain") {
    randomInitialState.startOfYearGrain should be >= (0)
    for (i <- 0 to 10) {
      deterministicStates(i).startOfYearGrain should be >= (0)
      if (i > 0) deterministicStates(i).startOfYearGrain should be (deterministicStates(i - 1).endOfYearGrain)
    }
  }

  test("harvestPerAcre") {
    randomInitialState.harvestPerAcre should (be >= (1) and be <= (5))
    for (i <- 0 to 10) {
      deterministicStates(i).harvestPerAcre should (be >= (1) and be <= (5))
    }
  }

  test("percentEatenByRats") {
    randomInitialState.percentEatenByRats should ( (be (0) or be >= (10)) and be <= (30) )
    for (i <- 0 to 10) {
      deterministicStates(i).percentEatenByRats should ( (be (0) or be >= (10)) and be <= (30) )
    }
  }

  test("plagueHappened") {
    randomInitialState.plagueHappened should be (false)
    deterministicStates.count(_.plagueHappened == false) should be > (0)
    deterministicStates.count(_.plagueHappened == true) should be > (0)
  }
  
  test("upcomingPricePerAcre") {
    randomInitialState.upcomingPricePerAcre should (be >= (17) and be <= (26))
    for (i <- 0 to 10) {
      deterministicStates(i).upcomingPricePerAcre should (be >= (17) and be <= (26))
      if (i > 0) deterministicStates(i).pricePerAcre should be (deterministicStates(i - 1).upcomingPricePerAcre)
    }
  }
  
  test("acresBought") {
    randomInitialState.acresBought should be >= (0)
    for (i <- 0 to 10) {
      deterministicStates(i).acresBought should be >= (0)
    }
  }

  test("acresSold") {
    val allSold = new KingdomState(deterministicStates(1), acresToBuyOrSell = -10000, grainForFood = 2000, acresToPlant = 200)
    allSold.acresSold should be (allSold.startOfYearAcreage)
        
    randomInitialState.acresSold should (be >= (0) and be <= (randomInitialState.startOfYearAcreage))
    for (i <- 0 to 10) {
      deterministicStates(i).acresSold should (be >= (0) and be <= (deterministicStates(i).startOfYearAcreage))
    }
  }

  test("grainUsedToBuyLand") {
    val allGrain = new KingdomState(deterministicStates(1), acresToBuyOrSell = 10000, grainForFood = 10000, acresToPlant = 10000)
    allGrain.grainUsedToBuyLand should be (allGrain.startOfYearGrain)
        
    randomInitialState.grainUsedToBuyLand should (be >= (0) and be <= (randomInitialState.startOfYearGrain))
    for (i <- 0 to 10) {
      deterministicStates(i).grainUsedToBuyLand should (be >= (0) and be <= (deterministicStates(i).startOfYearGrain))
    }
  }

  test("grainFromSaleOfLand") {
    val allAcres = new KingdomState(deterministicStates(1), acresToBuyOrSell = -10000, grainForFood = 10000, acresToPlant = 10000)
    allAcres.grainFromSaleOfLand should be (allAcres.startOfYearAcreage * allAcres.pricePerAcre)
        
    randomInitialState.grainFromSaleOfLand should (be >= (0) and be <= (randomInitialState.startOfYearAcreage * randomInitialState.pricePerAcre))
    for (i <- 0 to 10) {
      deterministicStates(i).grainFromSaleOfLand should (be >= (0) and be <= (deterministicStates(i).startOfYearAcreage * deterministicStates(i).pricePerAcre))
    }
  }

  test("grainAfterBartering") {
    val allGrain = new KingdomState(deterministicStates(1), acresToBuyOrSell = 10000, grainForFood = 10000, acresToPlant = 10000)
    allGrain.grainAfterBartering should be (0)
        
    randomInitialState.grainAfterBartering should be >= (0)
    for (i <- 0 to 10) {
      deterministicStates(i).grainAfterBartering should be >= (0)
    }
  }

  test("peopleFed") {
    val allGrain = new KingdomState(deterministicStates(1), acresToBuyOrSell = 0, grainForFood = 10000, acresToPlant = 10000)
//    allGrain.peopleFed should be (allGrain.startOfYearPopulation)

    randomInitialState.peopleFed should (be >= (0) and be <= (randomInitialState.startOfYearPopulation))
    for (i <- 0 to 10) {
      deterministicStates(i).peopleFed should (be >= (0) and be <= (deterministicStates(i).startOfYearPopulation))
    }
  }

  test("grainFedToPeople") {
    val allGrain = new KingdomState(deterministicStates(1), acresToBuyOrSell = -1 * deterministicStates(1).startOfYearAcreage, grainForFood = 100000, acresToPlant = 10000)
//    println("pop=%d gpp=%d gab=%d".format(allGrain.startOfYearPopulation, KingdomState.grainPerPerson, allGrain.grainAfterBartering))
    allGrain.grainFedToPeople should be (allGrain.startOfYearPopulation * KingdomState.grainPerPerson)

    randomInitialState.grainFedToPeople should (be >= (0) and be <= (randomInitialState.grainAfterBartering))
    for (i <- 0 to 10) {
      deterministicStates(i).grainFedToPeople should (be >= (0) and be <= (deterministicStates(i).grainAfterBartering))
    }
  }

  test("grainAfterFeeding") {
    val allGrain = new KingdomState(deterministicStates(1), acresToBuyOrSell = 0, grainForFood = 10000, acresToPlant = 10000)
    allGrain.grainAfterFeeding should be (allGrain.startOfYearGrain - allGrain.startOfYearPopulation * KingdomState.grainPerPerson)
        
    randomInitialState.grainAfterFeeding should (be >= (0) and be <= (randomInitialState.grainAfterBartering))
    for (i <- 0 to 10) {
      deterministicStates(i).grainAfterFeeding should (be >= (0) and be <= (deterministicStates(i).grainAfterBartering))
    }
  }

  test("plantingAcres") {
    val allPlanted = new KingdomState(deterministicStates(1), acresToBuyOrSell = 0, grainForFood = 0, acresToPlant = 100000)
    allPlanted.plantingAcres should be (allPlanted.endOfYearAcreage)
        
    randomInitialState.plantingAcres should (be >= (0) and be <= (randomInitialState.startOfYearAcreage + randomInitialState.acresBought - randomInitialState.acresSold))
    for (i <- 0 to 10) {
      deterministicStates(i).plantingAcres should (be >= (0) and be <= (deterministicStates(i).startOfYearAcreage + deterministicStates(i).acresBought - deterministicStates(i).acresSold))
    }
  }

  test("grainPlanted") {
    val allPlanted = new KingdomState(deterministicStates(1), acresToBuyOrSell = 0, grainForFood = 0, acresToPlant = 100000)
    allPlanted.grainPlanted should be (allPlanted.endOfYearAcreage / KingdomState.acresPerBushel)
        
    randomInitialState.grainPlanted should (be >= (0) and be <= (randomInitialState.grainAfterFeeding))
    for (i <- 0 to 10) {
      deterministicStates(i).grainPlanted should (be >= (0) and be <= (deterministicStates(i).grainAfterFeeding))
    }
  }

  test("acresPlanted") {
    val allPlanted = new KingdomState(deterministicStates(1), acresToBuyOrSell = 0, grainForFood = 0, acresToPlant = 100000)
//    println("pop=%d gpp=%d gab=%d".format(allGrain.startOfYearPopulation, KingdomState.grainPerPerson, allGrain.grainAfterBartering))
//    allPlanted.acresPlanted should be (allPlanted.startOfYearAcreage)
        
    randomInitialState.acresPlanted should (be >= (0) and be <= (randomInitialState.plantingAcres))
    for (i <- 0 to 10) {
      deterministicStates(i).acresPlanted should (be >= (0) and be <= (deterministicStates(i).plantingAcres))
    }
  }

  test("grainAfterPlanting") {
    val allPlanted = new KingdomState(deterministicStates(1), acresToBuyOrSell = 0, grainForFood = 0, acresToPlant = 100000)
    allPlanted.grainAfterPlanting should be (allPlanted.startOfYearGrain - allPlanted.endOfYearAcreage / KingdomState.acresPerBushel)
        
    randomInitialState.grainAfterPlanting should (be >= (0) and be <= (randomInitialState.grainAfterFeeding))
    for (i <- 0 to 10) {
      deterministicStates(i).grainAfterPlanting should (be >= (0) and be <= (deterministicStates(i).grainAfterFeeding))
    }
  }

  test("grainHarvested") {
    val allPlanted = new KingdomState(deterministicStates(1), acresToBuyOrSell = 0, grainForFood = 0, acresToPlant = 100000)
    allPlanted.grainHarvested should be (allPlanted.acresPlanted * allPlanted.harvestPerAcre)
        
    randomInitialState.grainHarvested should be (randomInitialState.acresPlanted * randomInitialState.harvestPerAcre)
    for (i <- 0 to 10) {
      deterministicStates(i).grainHarvested should be (deterministicStates(i).acresPlanted * deterministicStates(i).harvestPerAcre)
    }
  }

  test("grainAfterHarvest") {
    randomInitialState.grainAfterHarvest should be >= (randomInitialState.grainAfterPlanting)
    for (i <- 0 to 10) {
      deterministicStates(i).grainAfterHarvest should be >= (deterministicStates(i).grainAfterPlanting)
    }
  }

  test("grainEatenByRats") {
    if (randomInitialState.percentEatenByRats > 0)
      randomInitialState.grainEatenByRats should be <= (randomInitialState.grainAfterHarvest / 2)
    else
      randomInitialState.grainEatenByRats should be (0)
      
    for (i <- 0 to 10) {
      if (deterministicStates(i).percentEatenByRats > 0)
        deterministicStates(i).grainEatenByRats should be <= (deterministicStates(i).grainAfterHarvest / 2)
      else
        deterministicStates(i).grainEatenByRats should be (0)
    }
  }

  test("plagueVictims") {
    if (randomInitialState.plagueHappened)
      randomInitialState.plagueVictims should (be >= (0) and be <= (randomInitialState.startOfYearPopulation))
    else
      randomInitialState.plagueVictims should be (0)
      
    for (i <- 0 to 10) {
      if (deterministicStates(i).plagueHappened)
        deterministicStates(i).plagueVictims should (be >= (0) and be <= (deterministicStates(i).startOfYearPopulation))
      else
        deterministicStates(i).plagueVictims should be (0)
    }
  }

  test("postPlaguePopulation") {
    if (randomInitialState.plagueHappened)
      randomInitialState.postPlaguePopulation should (be >= (0) and be <= (randomInitialState.startOfYearPopulation))
    else
      randomInitialState.postPlaguePopulation should be (randomInitialState.startOfYearPopulation)
      
    for (i <- 0 to 10) {
      if (deterministicStates(i).plagueHappened)
        deterministicStates(i).postPlaguePopulation should (be >= (0) and be <= (deterministicStates(i).startOfYearPopulation))
      else
        deterministicStates(i).postPlaguePopulation should be (deterministicStates(i).startOfYearPopulation)
    }
  }

  test("starvationVictims") {
    randomInitialState.starvationVictims should (be >= (0) and be <= (randomInitialState.startOfYearPopulation))
      
    for (i <- 0 to 10) {
      deterministicStates(i).starvationVictims should (be >= (0) and be <= (deterministicStates(i).startOfYearPopulation))
    }
  }

  test("immigrants") {
    if (randomInitialState.starvationVictims == 0)
      randomInitialState.immigrants should (be >= (0) and be <= (randomInitialState.startOfYearPopulation))
    else
      randomInitialState.immigrants should be (0)
      
    for (i <- 0 to 10) {
      if (deterministicStates(i).starvationVictims == 0)
        deterministicStates(i).immigrants should (be >= (0) and be <= (deterministicStates(i).startOfYearPopulation))
      else
        deterministicStates(i).immigrants should be (0)
    }
  }

  test("endOfYearPopulation") {
    randomInitialState.endOfYearPopulation should be (100)
      
    for (i <- 0 to 10) {
      deterministicStates(i).endOfYearPopulation should be (expectedEOYPopulation(i))
    }
  }

  test("endOfYearAcreage") {
    randomInitialState.endOfYearAcreage should be (1000)
      
    for (i <- 0 to 10) {
      deterministicStates(i).endOfYearAcreage should be (expectedEOYAcreage(i))
    }
  }

  test("endOfYearGrain") {
    randomInitialState.endOfYearGrain should be (2800)
      
    for (i <- 0 to 10) {
      deterministicStates(i).endOfYearGrain should be (expectedEOYGrain(i))
    }
  }
    
  test("stillInOffice") {
    randomInitialState.stillInOffice should be (true)

    for (i <- 0 to 10) {
      deterministicStates(i).stillInOffice should be (expectedEOYStillInOffice(i))
    }
  }
}


class KingdomStateCompanionObjectTest extends FunSuite with ShouldMatchers {
  
  test("price per acre in expected range") {
    val list = for (i <- 1 to 100) yield KingdomState.randomPricePerAcre(i == 1)

    for (v <- list) {
      v should be >= 17
      v should be <= 26
    }
  }

  test("yield per acre in expected range") {
    val list = for (i <- 1 to 100) yield KingdomState.randomYieldPerAcre(i == 1)

    for (v <- list) {
      v should be >= 1
      v should be <= 5
    }
  }

  test("rat percentage in expected range") {
    val list = for (i <- 1 to 1000) yield KingdomState.randomRatPercent(i == 1)

    list.head should be >= (10) // deterministic case should have some rats
    
    for (v <- list) {
      v should (be === (0) or be >= (10))
      v should be <= 30
    }
    
    list.count(_ == 0) should be > (list.length * 51 / 100) // no rats most of the time
    list.count(_ > 0) should be > (list.length * 10 / 100) // some rats some of the time
  }
  
  test("plague occurrence in expected range") {
    val list = for (i <- 1 to 1000) yield KingdomState.randomPlagueHappened(i == 1)

    list.head should be >= (true) // deterministic case should have a plague
    list.count(_ == true) should be > (list.length * 10 / 100) // random plague at least y% of time
    list.count(_ == false) should be > (list.length * 51 / 100) // no plague most of time
  }
}
*/
