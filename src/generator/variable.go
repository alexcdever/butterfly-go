package generator

import "time"

var (
	instanceList []*ButterflyList
	// period defines the  execution period of the checkUnusedIDListCount in the gorountine
	period = 300 * time.Second
)

func init() {
	go checkUnusedIDListCount()
}

func checkUnusedIDListCount() {
	for {
		for _, instance := range instanceList {
			instance.mutex.Lock()
			defer instance.mutex.Unlock()

			if instance.AtLeastCount < len(instance.UnusedIDList) {
				instance.UnusedIDList = append(instance.UnusedIDList, instance.generator.GenerateInBatches(instance.IncreaseCount)...)
			}

		}
		time.Sleep(period)
	}
}
