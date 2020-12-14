package igbinary

func (d *Decoder) uint8() (uint8, error) {
	c, err := d.readCode()
	if err != nil {
		return 0, err
	}
	return c, nil
}

func (d *Decoder) uint16() (uint16, error) {
	b, err := d.readN(2)
	if err != nil {
		return 0, err
	}
	return (uint16(b[0]) << 8) | uint16(b[1]), nil
}

func (d *Decoder) uint32() (uint32, error) {
	b, err := d.readN(4)
	if err != nil {
		return 0, err
	}
	n := (uint32(b[0]) << 24) |
		(uint32(b[1]) << 16) |
		(uint32(b[2]) << 8) |
		uint32(b[3])
	return n, nil
}
