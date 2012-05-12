package kingdomstate

import (
	"math/rand"
	"fmt"
)

const (
	GrainPerPerson = 20
	AcresPerBushel = 2
	AcresPerPerson = 20
)

func RandomPricePerAcre(randgen *rand.Rand) uint {
	if randgen == nil {
		return 21
	}
	return uint(randgen.Intn(10) + 17)
}

func RandomYieldPerAcre(randgen *rand.Rand) uint {
	if randgen == nil {
		return 3
	}
	return uint(randgen.Intn(5) + 1)
}

func RandomRatPercent(randgen *rand.Rand) uint {
	if randgen == nil {
		return 10
	}
	if randgen.Float32() < 0.40 {
		return uint(randgen.Intn(21) + 10)
	}
	return 0
}

func RandomPlagueHappened(randgen *rand.Rand, year uint) bool {
	if randgen == nil {
		return (year%4 == 0)
	}
	return (randgen.Float32() > 0.85)
}

func min(i, j uint) uint {
	if i < j { return i }
	return j
}

func max(i, j uint) uint {
	if i > j { return i }
	return j
}

type KingdomState struct {
	randgen *rand.Rand
	
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
	
	starvationVictims uint
	plagueVictims uint
	immigrants uint
	grainHarvested uint
	grainEatenByRats uint
}

func (ks *KingdomState) SetupInitialState(randgen *rand.Rand) {
	ks.randgen = randgen
	
	ks.stillInOffice = true
	
	ks.population = 100
	ks.acreage = 1000
	ks.grain = 2800
	
	ks.harvestPerAcre = 3
	ks.percentEatenByRats = 10
	ks.nextYearPricePerAcre = RandomPricePerAcre(ks.randgen)

	ks.starvationVictims = 0
	ks.plagueVictims = 0
	ks.immigrants = 5
	ks.grainHarvested = 3000
	ks.grainEatenByRats = 400
}

func (ks *KingdomState) TallyUpYear(acresToBuy, acresToSell, grainForFood, acresToPlant uint) {
	if !ks.stillInOffice { panic(0) }
	
	var startOfYearPopulation uint = ks.population
	
	// Cummulative values
	ks.yearOfRule++
	ks.pricePerAcre = ks.nextYearPricePerAcre

	// Random events
	ks.harvestPerAcre = RandomYieldPerAcre(ks.randgen)
	ks.percentEatenByRats = RandomRatPercent(ks.randgen)
	ks.plagueHappened = RandomPlagueHappened(ks.randgen, ks.yearOfRule)  
	ks.nextYearPricePerAcre = RandomPricePerAcre(ks.randgen)
	
	// Buy land
	grainUsedToBuyLand := min(acresToBuy * ks.pricePerAcre, ks.grain)
	ks.grain -= grainUsedToBuyLand
	ks.acreage += grainUsedToBuyLand / ks.pricePerAcre

	// Sell land
	grainFromSaleOfLand := min(acresToSell, ks.acreage) * ks.pricePerAcre
	ks.grain += grainFromSaleOfLand
	ks.acreage -= grainFromSaleOfLand / ks.pricePerAcre
	
	// Feed the people
	peopleFed := min( min(ks.grain, grainForFood) / GrainPerPerson, ks.population)
	ks.grain -= peopleFed * GrainPerPerson
	
	// Plant the fields
	acresForPlanting :=	min( min(acresToPlant, ks.population * AcresPerPerson), ks.acreage )
	grainPlanted := min(ks.grain, acresForPlanting / AcresPerBushel)
	acresPlanted := grainPlanted * AcresPerBushel
	ks.grain -= grainPlanted
	grainAfterPlanting := ks.grain
	
	// Harvest grain and deal with the rats
	ks.grainHarvested = acresPlanted * ks.harvestPerAcre
	ks.grain += ks.grainHarvested
	ks.grainEatenByRats = ks.percentEatenByRats * ks.grain / 100
	ks.grain -= ks.grainEatenByRats
	
	// Adjust population counts
	if ks.plagueHappened {
		ks.plagueVictims = ks.population / 2
		ks.population -= ks.plagueVictims
	} else {
		ks.plagueVictims = 0
	}
	
	if ks.population > peopleFed {
		// Starvation occurs if not everyone was fed
		ks.starvationVictims = ks.population - peopleFed
		ks.population -= ks.starvationVictims
	} else {
		ks.starvationVictims = 0
	}
	
	if ks.population > 0 && ks.starvationVictims == 0 {
		// Allow immigrants if nobody starved and there are still people around
		ks.immigrants = (20 * ks.acreage + grainAfterPlanting) / (100 * ks.population) + 1
		ks.population += ks.immigrants
	} else {
		ks.immigrants = 0
	}

	// Determine if the game is over
	ks.stillInOffice = (
		ks.yearOfRule < 10 &&
		ks.population > 0 &&
		ks.starvationVictims < 45 * startOfYearPopulation / 100)
}

func (ks KingdomState) YearOfRule() uint {
	return ks.yearOfRule
}
func (ks KingdomState) Population() uint {
	return ks.population
}
func (ks KingdomState) Grain() uint {
	return ks.grain
}
func (ks KingdomState) Acreage() uint {
	return ks.acreage
}

func (ks KingdomState) StillInOffice() bool {
	return ks.stillInOffice
}

func (ks KingdomState) PrintSummary() {
	fmt.Printf("___________________________________________________________________")
	fmt.Printf("\nO Great Hammurabi!\n")
	fmt.Printf("You are in year %d of your ten year rule.\n", ks.yearOfRule + 1)
	if (ks.plagueVictims > 0) {
		fmt.Printf("A horrible plague killed %d people.\n", ks.plagueVictims)
	}
	fmt.Printf("In the previous year %d people starved to death.\n", ks.starvationVictims)
	fmt.Printf("In the previous year %d people entered the kingdom.\n", ks.immigrants)
	fmt.Printf("The population is now %d.\n", ks.population)
	fmt.Printf("We harvested %d bushels at %d bushels per acre.\n", ks.grainHarvested, ks.harvestPerAcre)
	if (ks.grainEatenByRats > 0) {
		fmt.Printf("*** Rats destroyed %d bushels, leaving %d bushels in storage.\n", ks.grainEatenByRats, ks.grain)
	} else {
		fmt.Printf("We have %d bushels of grain in storage.\n", ks.grain)
	}
	fmt.Printf("The city owns %d acres of land.\n", ks.acreage)
	fmt.Printf("Land is currently worth %d bushels per acre.\n", ks.nextYearPricePerAcre)
}

