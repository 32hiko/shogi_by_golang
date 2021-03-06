package shogi

type TJoseki struct {
	FixOpening map[int]*TMove
	SFENMap    map[string]*TMove
}

func NewJoseki() *TJoseki {
	fix_opening := CreateFixOpening()
	sfen_map := CreateSFENMap()
	joseki := TJoseki{
		FixOpening: fix_opening,
		SFENMap:    sfen_map,
	}
	return &joseki
}

func CreateFixOpening() map[int]*TMove {
	m := make(map[int]*TMove)
	// 先手の場合、初手78飛
	// moveを作るために、ダミーの駒が必要。（直したい）
	// m[1] = NewMove(NewKoma(1, Hi, 2, 8, Sente), TPosition(complex(7, 8)), 0)
	// 後手の場合、初手32飛→26歩の場合即終了。だめ。
	// m[2] = NewMove(NewKoma(2, Hi, 8, 2, Gote), TPosition(complex(3, 2)), 0)
	// 角頭を受けるために、角道を開ける
	// m[2] = NewMove(NewKoma(2, Fu, 3, 3, Gote), TPosition(complex(3, 4)), 0)
	// m[4] = NewMove(NewKoma(4, Kaku, 2, 2, Gote), TPosition(complex(3, 3)), 0)
	return m
}

func CreateSFENMap() map[string]*TMove {
	m := make(map[string]*TMove)
	// 先手の場合
	{
		// ▲7六歩
		m["lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL/ b -"] = NewMove(NewKoma(1, Fu, 7, 7, Sente), TPosition(complex(7, 6)), 0)
		// -> △3四歩 -> ▲7五歩
		m["lnsgkgsnl/1r5b1/pppppp1pp/6p2/9/2P6/PP1PPPPPP/1B5R1/LNSGKGSNL/ b -"] = NewMove(NewKoma(3, Fu, 7, 6, Sente), TPosition(complex(7, 5)), 0)
		// -> △8四歩 -> ▲7八飛
		m["lnsgkgsnl/1r5b1/p1pppp1pp/1p4p2/2P6/9/PP1PPPPPP/1B5R1/LNSGKGSNL/ b -"] = NewMove(NewKoma(5, Hi, 2, 8, Sente), TPosition(complex(7, 8)), 0)
		// -> △8五歩 -> ▲7六飛
		m["lnsgkgsnl/1r5b1/p1pppp1pp/6p2/1pP6/9/PP1PPPPPP/1BR6/LNSGKGSNL/ b -"] = NewMove(NewKoma(7, Hi, 7, 8, Sente), TPosition(complex(7, 6)), 0)
		// ▲7六歩 -> △8四歩 -> ▲7八飛
		m["lnsgkgsnl/1r5b1/p1ppppppp/1p7/9/2P6/PP1PPPPPP/1B5R1/LNSGKGSNL/ b -"] = NewMove(NewKoma(3, Hi, 2, 8, Sente), TPosition(complex(7, 8)), 0)
		// -> △8五歩 -> ▲7七角
		m["lnsgkgsnl/1r5b1/p1ppppppp/9/1p7/2P6/PP1PPPPPP/1BR6/LNSGKGSNL/ b -"] = NewMove(NewKoma(5, Kaku, 8, 8, Sente), TPosition(complex(7, 7)), 0)
		// ▲7六歩 -> △3四歩 -> ▲7五歩 -> △6二銀 -> ▲7八飛
		m["ln1gkgsnl/1r1s3b1/pppppp1pp/6p2/2P6/9/PP1PPPPPP/1B5R1/LNSGKGSNL/ b -"] = NewMove(NewKoma(5, Hi, 2, 8, Sente), TPosition(complex(7, 8)), 0)
		// ▲7六歩 -> △3二金 -> ▲7八飛
		m["lnsgk1snl/1r4gb1/ppppppppp/9/9/2P6/PP1PPPPPP/1B5R1/LNSGKGSNL/ b -"] = NewMove(NewKoma(3, Hi, 2, 8, Sente), TPosition(complex(7, 8)), 0)
		// ▲7六歩 -> △3四歩 -> ▲7五歩 -> △8八角成 -> ▲同銀
		m["lnsgkgsnl/1r7/pppppp1pp/6p2/2P6/9/PP1PPPPPP/1+b5R1/LNSGKGSNL/ b b"] = NewMove(NewKoma(5, Gin, 7, 9, Sente), TPosition(complex(8, 8)), 0)
		// -> △4五角 -> ▲7六角
		m["lnsgkgsnl/1r7/pppppp1pp/6p2/2P2b3/9/PP1PPPPPP/1S5R1/LN1GKGSNL/ b B"] = NewMove(NewKoma(7, Kaku, 0, 0, Sente), TPosition(complex(7, 6)), 0)
		// ▲7六歩 -> △3四歩 -> ▲7五歩 -> △4二玉 -> ▲7八飛
		m["lnsg1gsnl/1r3k1b1/pppppp1pp/6p2/2P6/9/PP1PPPPPP/1B5R1/LNSGKGSNL/ b -"] = NewMove(NewKoma(5, Hi, 2, 8, Sente), TPosition(complex(7, 8)), 0)
		// -> △6二銀 -> ▲6六歩
		m["ln1g1gsnl/1r1s1k1b1/pppppp1pp/6p2/2P6/9/PP1PPPPPP/1BR6/LNSGKGSNL/ b -"] = NewMove(NewKoma(7, Fu, 6, 7, Sente), TPosition(complex(6, 6)), 0)
		// ▲7六歩 -> △8四歩 -> ▲7八飛 -> △8五歩 -> ▲7七角 -> △3四歩 -> ▲6六歩
		m["lnsgkgsnl/1r5b1/p1pppp1pp/6p2/1p7/2P6/PPBPPPPPP/2R6/LNSGKGSNL/ b -"] = NewMove(NewKoma(7, Fu, 6, 7, Sente), TPosition(complex(6, 6)), 0)
	}
	// 後手の場合
	{
		// ▲2六歩 -> △3四歩
		m["lnsgkgsnl/1r5b1/ppppppppp/9/9/7P1/PPPPPPP1P/1B5R1/LNSGKGSNL/ w -"] = NewMove(NewKoma(2, Fu, 3, 3, Gote), TPosition(complex(3, 4)), 0)
		// -> ▲2五歩 -> △3三角
		m["lnsgkgsnl/1r5b1/pppppp1pp/6p2/7P1/9/PPPPPPP1P/1B5R1/LNSGKGSNL/ w -"] = NewMove(NewKoma(4, Kaku, 2, 2, Gote), TPosition(complex(3, 3)), 0)
		// -> ▲7六歩 -> △4四歩
		m["lnsgkgsnl/1r5b1/pppppp1pp/6p2/9/2P4P1/PP1PPPP1P/1B5R1/LNSGKGSNL/ w -"] = NewMove(NewKoma(4, Fu, 4, 3, Gote), TPosition(complex(4, 4)), 0)
		// (合流) ▲2五歩 -> △3三角
		m["lnsgkgsnl/1r5b1/ppppp2pp/5pp2/7P1/2P6/PP1PPPP1P/1B5R1/LNSGKGSNL/ w -"] = NewMove(NewKoma(6, Kaku, 2, 2, Gote), TPosition(complex(3, 3)), 0)
		// (合流) ▲7六歩 -> △4四歩
		m["lnsgkgsnl/1r7/ppppppbpp/6p2/7P1/2P6/PP1PPPP1P/1B5R1/LNSGKGSNL/ w -"] = NewMove(NewKoma(6, Fu, 4, 3, Gote), TPosition(complex(4, 4)), 0)
		// (合流) ▲6八玉 -> △3二金
		m["lnsgkgsnl/1r5b1/ppppp2pp/5pp2/9/2P4P1/PP1PPPP1P/1B1K3R1/LNSG1GSNL/ w -"] = NewMove(NewKoma(6, Kin, 4, 1, Gote), TPosition(complex(3, 2)), 0)
		// ▲4八銀 -> △3二飛
		m["lnsgkgsnl/1r7/ppppp1bpp/5pp2/7P1/2P6/PP1PPPP1P/1B3S1R1/LNSGKG1NL/ w -"] = NewMove(NewKoma(8, Hi, 8, 2, Gote), TPosition(complex(3, 2)), 0)
	}
	{
		// ▲7六歩 -> △3四歩
		m["lnsgkgsnl/1r5b1/ppppppppp/9/9/2P6/PP1PPPPPP/1B5R1/LNSGKGSNL/ w -"] = NewMove(NewKoma(2, Fu, 3, 3, Gote), TPosition(complex(3, 4)), 0)
		// ▲2六歩 -> △4四歩は合流
		// ▲7六歩 -> △3四歩 -> ▲4八銀 -> △3二飛
		m["lnsgkgsnl/1r5b1/pppppp1pp/6p2/9/2P6/PP1PPPPPP/1B3S1R1/LNSGKG1NL/ w -"] = NewMove(NewKoma(4, Hi, 8, 2, Gote), TPosition(complex(3, 2)), 0)
		// ▲7六歩 -> △3四歩 -> ▲4八銀 -> △3二飛 -> ▲6八玉 -> △4四歩
		m["lnsgkgsnl/6rb1/pppppp1pp/6p2/9/2P6/PP1PPPPPP/1B1K1S1R1/LNSG1G1NL/ w -"] = NewMove(NewKoma(6, Fu, 4, 3, Gote), TPosition(complex(4, 4)), 0)
	}
	return m
}
