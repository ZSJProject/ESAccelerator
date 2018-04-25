package Core

import (
	"log"
	"reflect"
	"sync"
	"time"
)

const ESQueuePopFBoxLimit = 1000

type ESTimestamp int64
type ESQueueReturnType []*ESRequest

type ESQueueItem struct {
	Item      ESQueueReturnType
	Timestamp ESTimestamp
}

type Queue struct {
	Items  []ESQueueItem
	Latest ESTimestamp

	Mutex sync.Mutex
}

func (this *Queue) Length() int {
	return len(this.Items)
}

func (this *Queue) Push(Item interface{}, Manual *ESTimestamp) ESTimestamp {
	Target := func() ESQueueReturnType {
		switch Item.(type) {
		case *ESRequest:
			return ESQueueReturnType{Item.(*ESRequest)}

		case ESQueueReturnType:
			return Item.(ESQueueReturnType)
		}

		log.Fatalf("Push 메소드는 0번째 인자에 다음과 같은 형식을 허용하지 않습니다: %s", reflect.TypeOf(Item).Name())

		return nil
	}()

	CurrentTime := func() ESTimestamp {
		if Manual != nil {
			return ESTimestamp(*Manual)

		}

		return ESTimestamp(time.Now().UnixNano())
	}()

	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	this.Items = append(this.Items, ESQueueItem{Target, CurrentTime})

	if this.Latest == -1 || this.Latest < CurrentTime {
		this.Latest = CurrentTime
	}

	return this.Latest
}

func (this *Queue) Pop() *ESQueueReturnType {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	return this.UnsafePop()
}

func (this *Queue) UnsafePop() *ESQueueReturnType {
	if this.Length() > 0 {
		Target := this.Items[this.Length()-1]

		this.Items = this.Items[:this.Length()-1]

		if this.Length() > 0 {
			this.Latest = this.Items[this.Length()-1].Timestamp
		}

		return &Target.Item
	}

	return nil
}

func (this *Queue) MLength(Td ESTimestamp) int {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	return this.UnsafeMLength(Td)
}

func (this *Queue) UnsafeMLength(Threshold ESTimestamp) int {
	if Threshold <= 0 {
		return this.Length()
	}

	if Threshold < this.Latest && Threshold != -1 {
		Seek := this.Length()
		Range := this.Latest - Threshold

		for {
			if Seek == 0 || this.Items[Seek-1].Timestamp < Range {
				break
			}

			Seek--
		}

		return this.Length() - Seek
	}

	return 0
}

func (this *Queue) MPush(Threshold interface{}, Items ...*ESRequest) ESTimestamp {
	switch Threshold.(type) {
	case ESTimestamp:
		{
			if Threshold.(ESTimestamp) <= 0 {
				log.Fatalf("MPush 메소드는 0번째 인자의 0이하의 값을 허용하지 않습니다: %d", Threshold.(ESTimestamp))
			}
		}

	default:
		{
			if Threshold != nil {
				log.Fatal("MPush 메소드는 ESTimestamp 형식이 아닌 0번째 인자에 어떠한 값도 허용하지 않습니다.")
			}
		}
	}

	Range := func() ESTimestamp {
		if Threshold == nil {
			return 0
		}

		return Threshold.(ESTimestamp)
	}()

	Packed := make(ESQueueReturnType, len(Items))

	for _, V := range Items {
		Packed = append(Packed, V)
	}

	return this.Push(Packed, &Range)
}

func (this *Queue) MPop(Threshold ESTimestamp) *ESQueueReturnType {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	Range := this.UnsafeMLength(Threshold)

	if Range != 0 {
		Items := func(Items ...ESQueueItem) *[]ESQueueReturnType {
			Outgoing := make([]ESQueueReturnType, Range)

			for Acc, V := range Items {
				Outgoing[Acc] = V.Item
			}

			return &Outgoing

		}(this.Items[this.Length()-Range : this.Length()]...)

		FBox := make(ESQueueReturnType, 0, func(ShouldMake int) int {
			if ShouldMake > ESQueuePopFBoxLimit {
				return ESQueuePopFBoxLimit
			}

			return ShouldMake

		}(this.Length()*2))

		Product := make(ESQueueReturnType, 0, func() int {
			Sum := 0

			for _, V := range *Items {
				for _, V_ := range V {
					FBox = append(FBox, V_)
				}

				Sum += len(V)
			}

			return Sum
		}())

		Product = append(Product, FBox...)

		for Acc := 0; Acc < Range; Acc++ {
			this.UnsafePop()
		}

		return &Product
	}

	return nil
}

func CreateNewQueue() Queue {
	return Queue{Latest: -1}
}
