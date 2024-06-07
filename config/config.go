package config

import "time"

var (
	ValidatorIndx = []uint32{
		216529, 216530, 216531, 216536, 216537, 216538, 216539, 216540, 216541, 216542, 216543,
		216544, 216545, 216549, 235953, 254984, 268039, 268323, 269327, 269437, 273801, 273875,
		275464, 275768, 278948, 285161, 286076, 286247, 301641, 301642, 301766, 302947, 314637,
		321664, 332391, 335915, 337944, 339070, 346309, 351564, 353018, 353631, 356336, 372937,
		373015, 377457, 379082, 385084, 386549, 386673, 387901, 392258, 393239, 393400, 394710,
		394711, 399322, 415326, 415353, 431732, 438657, 452943, 453472, 457331, 460097, 464035,
		465083, 470270, 471464, 472836, 473956, 474004, 474125, 474129, 478239, 478243, 478244,
		478245, 479138, 484156, 484574, 484577, 485001, 486823, 487030, 488161, 495953, 495954,
		502083, 503169, 503170, 503174, 503281, 510557, 512445, 512648, 512655, 513067, 513181,
		513631}
	Port                     = "8080"
	RateLimitPerMinute uint8 = 10
	PollInterval             = time.Second * 30
)
