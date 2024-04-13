package kvstore

import (
	"os"
)

// BlockDeviceAdapter 提供了对块设备的基本操作
type BlockDeviceAdapter struct {
}

// NewBlockDeviceAdapter 创建一个新的BlockDeviceAdapter实例
func NewBlockDeviceAdapter() *BlockDeviceAdapter {
	return &BlockDeviceAdapter{}
}

// ReadBlock 从块设备读取数据
func (bda *BlockDeviceAdapter) ReadBlock(offset, size int64) ([]byte, error) {
	// 打开设备文件
	file, err := os.Open("")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 设置读取的偏移量
	_, err = file.Seek(offset, 0)
	if err != nil {
		return nil, err
	}

	// 读取数据
	buf := make([]byte, size)
	_, err = file.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// WriteBlock 向块设备写入数据
func (bda *BlockDeviceAdapter) WriteBlock(offset int64, data []byte) error {
	// 以读写模式打开设备文件
	file, err := os.OpenFile("", os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	// 设置写入的偏移量
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// 写入数据
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
