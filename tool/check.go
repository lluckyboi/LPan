package tool

func LengthCheck(ss string) (string, bool) {
	if len(ss) > 20 || len(ss) < 2 {
		err := "超过长度限制！"
		return err, false
	}
	return "", true
}

func PostLengthCheck(ss string) (string, bool) {
	if len(ss) > 20 {
		err := "超过长度限制！"
		return err, false
	}
	return "", true
}

func ScoreCheck(score float32) (string, bool) {
	if score > 10 || score < 0 {
		err := "长度不符合要求"
		return err, false
	}
	return "", true
}
