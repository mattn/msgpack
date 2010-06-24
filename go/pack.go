package msgpack

import (
    "io"
    "os"
    "unsafe"
    "reflect"
);

func packUint8(writer io.Writer, value uint8) (n int, err os.Error) {
    if value < 128 {
        return writer.Write([]byte { byte(value) })
    }
    return writer.Write([]byte { 0xcc, byte(value) })
}

func packUint16(writer io.Writer, value uint16) (n int, err os.Error) {
    if value < 128 {
        return writer.Write([]byte { byte(value) })
    } else if value < 256 {
        return writer.Write([]byte { 0xcc, byte(value) })
    }
    return writer.Write([]byte { 0xcd, byte(value >> 8), byte(value) })
}

func packUint32(writer io.Writer, value uint32) (n int, err os.Error) {
    if value < 128 {
        return writer.Write([]byte { byte(value) })
    } else if value < 256 {
        return writer.Write([]byte { 0xcc, byte(value) })
    } else if value < 65536 {
        return writer.Write([]byte { 0xcd, byte(value >> 8), byte(value) })
    }
    return writer.Write([]byte { 0xce, byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value) })
}

func packUint64(writer io.Writer, value uint64) (n int, err os.Error) {
    if value < 128 {
        return writer.Write([]byte { byte(value) })
    } else if value < 256 {
        return writer.Write([]byte { 0xcc, byte(value) })
    } else if value < 65536 {
        return writer.Write([]byte { 0xcd, byte(value >> 8), byte(value) })
    } else if value < 4294967296 {
        return writer.Write([]byte { 0xce, byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value) })
    }
    return writer.Write([]byte { 0xcf, byte(value >> 56), byte(value >> 48), byte(value >> 40), byte(value >> 32), byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value) })
}

func PackUint(writer io.Writer, value *reflect.UintValue) (n int, err os.Error) {
    switch value.Type().Size() {
    case 1:
        return packUint8(writer, uint8(value.Get()))
    case 2:
        return packUint16(writer, uint16(value.Get()))
    case 4:
        return packUint32(writer, uint32(value.Get()))
    case 8:
        return packUint64(writer, uint64(value.Get()));
    }
    return 0, os.ENOENT // never get here
}


func packInt8(writer io.Writer, value int8) (n int, err os.Error) {
    if value < -32 {
        return writer.Write([]byte { 0xd0, byte(value) })
    }
    return writer.Write([]byte { byte(value) })
}

func packInt16(writer io.Writer, value int16) (n int, err os.Error) {
    if value < -128 || value >= 128 {
        return writer.Write([]byte { 0xd1, byte(uint16(value) >> 8), byte(value) })
    } else if value < -32 {
        return writer.Write([]byte { 0xd0, byte(value) })
    }
    return writer.Write([]byte { byte(value) })
}

func packInt32(writer io.Writer, value int32) (n int, err os.Error) {
    if value < -32768 || value >= 32768 {
        return writer.Write([]byte { 0xd2, byte(uint32(value) >> 24), byte(uint32(value) >> 16), byte(uint32(value) >> 8), byte(value) })
    } else if value < -128 {
        return writer.Write([]byte { 0xd1, byte(uint32(value) >> 8), byte(value) })
    } else if value < -32 {
        return writer.Write([]byte { 0xd0, byte(value) })
    } else if value >= 128 {
        return writer.Write([]byte { 0xd1, byte(uint32(value) >> 8), byte(value) })
    }
    return writer.Write([]byte { byte(value) })
}

func packInt64(writer io.Writer, value int64) (n int, err os.Error) {
    if value < -2147483648 || value >= 2147483648 {
        return writer.Write([]byte { 0xd3, byte(uint64(value) >> 56), byte(uint64(value) >> 48), byte(uint64(value) >> 40), byte(uint64(value) >> 32), byte(uint64(value) >> 24), byte(uint64(value) >> 16), byte(uint64(value) >> 8), byte(value) })
    } else if value < -32768 || value >= 32768 {
        return writer.Write([]byte { 0xd2, byte(uint64(value) >> 24), byte(uint64(value) >> 16), byte(uint64(value) >> 8), byte(value) })
    } else if value < -128 || value >= 128 {
        return writer.Write([]byte { 0xd1, byte(uint64(value) >> 8), byte(value) })
    } else if value < -32 {
        return writer.Write([]byte { 0xd0, byte(value) })
    }
    return writer.Write([]byte { byte(value) })
}

func PackInt(writer io.Writer, value *reflect.IntValue) (n int, err os.Error) {
    switch value.Type().Size() {
    case 1:
        return packInt8(writer, int8(value.Get()))
    case 2:
        return packInt16(writer, int16(value.Get()))
    case 4:
        return packInt32(writer, int32(value.Get()))
    case 8:
        return packInt64(writer, int64(value.Get()))
    }
    return 0, os.ENOENT // never get here
}

func PackNil(writer io.Writer) (n int, err os.Error) {
    return writer.Write([]byte{ 0xc0 })
}

func PackBool(writer io.Writer, value *reflect.BoolValue) (n int, err os.Error) {
    var code byte;
    if value.Get() {
        code = 0xc3
    } else {
        code = 0xc2
    }
    return writer.Write([]byte{ code })
}

func packFloat32(writer io.Writer, value float32) (n int, err os.Error) {
    return packUint32(writer, *(*uint32)(unsafe.Pointer(&value)))
}

func packFloat64(writer io.Writer, value float64) (n int, err os.Error) {
    return packUint64(writer, *(*uint64)(unsafe.Pointer(&value)))
}

func PackFloat(writer io.Writer, value *reflect.FloatValue) (n int, err os.Error) {
    switch value.Type().Size() {
    case 4:
        return packFloat32(writer, float32(value.Get()))
    case 8:
        return packFloat64(writer, float64(value.Get()));
    }
    return 0, os.ENOENT // never get here
}

func packBytes(writer io.Writer, value []byte) (n int, err os.Error) {
    if len(value) < 32 {
        n1, err := writer.Write([]byte { 0xa0 | uint8(len(value)) })
        if err != nil { return n1, err }
        n2, err := writer.Write(value)
        return n1 + n2, err
    } else if len(value) < 65536 {
        n1, err := writer.Write([]byte { 0xda, byte(len(value) >> 16), byte(len(value)) })
        if err != nil { return n1, err }
        n2, err := writer.Write(value)
        return n1 + n2, err
    }
    n1, err := writer.Write([]byte { 0xdb, byte(len(value) >> 24), byte(len(value) >> 16), byte(len(value) >> 8), byte(len(value)) })
    if err != nil { return n1, err }
    n2, err := writer.Write(value)
    return n1 + n2, err
}

func PackString(writer io.Writer, value *reflect.StringValue) (n int, err os.Error) {
    return packBytes(writer, ([]byte)(value.Get()));
}

func PackSlice(writer io.Writer, value *reflect.SliceValue) (n int, err os.Error) {
    l := value.Len()
    if l < 16 {
        n, err := writer.Write([]byte { 0x90 | byte(l) })
        if err != nil { return n, err }
        for i := 0; i < l; i++ {
            _n, err := PackValue(writer, value.Elem(i))
            if err != nil { return n, err }
            n += _n
        }
    } else if l < 65536 {
        n, err := writer.Write([]byte { 0xdc, byte(l >> 8), byte(l) })
        if err != nil { return n, err }
        for i := 0; i < l; i++ {
            _n, err := PackValue(writer, value.Elem(i))
            if err != nil { return n, err }
            n += _n
        }
    } else {
        n, err := writer.Write([]byte { 0xdd, byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l) })
        if err != nil { return n, err }
        for i := 0; i < l; i++ {
            _n, err := PackValue(writer, value.Elem(i))
            if err != nil { return n, err }
            n += _n
        }
    }
    return n, nil
}

func PackArray(writer io.Writer, value *reflect.ArrayValue) (n int, err os.Error) {
    l := value.Len()
    if l < 16 {
        n, err := writer.Write([]byte { 0x90 | byte(l) })
        if err != nil { return n, err }
        for i := 0; i < l; i++ {
            _n, err := PackValue(writer, value.Elem(i))
            if err != nil { return n, err }
            n += _n
        }
    } else if l < 65536 {
        n, err := writer.Write([]byte { 0xdc, byte(l >> 8), byte(l) })
        if err != nil { return n, err }
        for i := 0; i < l; i++ {
            _n, err := PackValue(writer, value.Elem(i))
            if err != nil { return n, err }
            n += _n
        }
    } else {
        n, err := writer.Write([]byte { 0xdd, byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l) })
        if err != nil { return n, err }
        for i := 0; i < l; i++ {
            _n, err := PackValue(writer, value.Elem(i))
            if err != nil { return n, err }
            n += _n
        }
    }
    return n, nil
}

func PackMap(writer io.Writer, value *reflect.MapValue) (n int, err os.Error) {
    keys := value.Keys()
    if value.Len() < 16 {
        n, err := writer.Write([]byte { 0x80 | byte(len(keys)) })
        if err != nil { return n, err }
        for _, k := range keys {
            _n, err := PackValue(writer, k)
            if err != nil { return n, err }
            n += _n
            _n, err = PackValue(writer, value.Elem(k))
            if err != nil { return n, err }
            n += _n
        }
    } else if value.Len() < 65536 {
        n, err := writer.Write([]byte { 0xde, byte(len(keys) >> 8), byte(len(keys)) })
        if err != nil { return n, err }
        for _, k := range keys {
            _n, err := PackValue(writer, k)
            if err != nil { return n, err }
            n += _n
            _n, err = PackValue(writer, value.Elem(k))
            if err != nil { return n, err }
            n += _n
        }
    } else {
        n, err := writer.Write([]byte { 0xdf, byte(len(keys) >> 24), byte(len(keys) >> 16), byte(len(keys) >> 8), byte(len(keys)) })
        if err != nil { return n, err }
        for _, k := range keys {
            _n, err := PackValue(writer, k)
            if err != nil { return n, err }
            n += _n
            _n, err = PackValue(writer, value.Elem(k))
            if err != nil { return n, err }
            n += _n
        }
    }
    return n, nil
}

func PackValue(writer io.Writer, value reflect.Value) (n int, err os.Error) {
    switch _value := value.(type) {
    case *reflect.BoolValue: return PackBool(writer, value.(*reflect.BoolValue));
    case *reflect.UintValue: return PackUint(writer, value.(*reflect.UintValue));
    case *reflect.IntValue: return PackInt(writer, value.(*reflect.IntValue));
    case *reflect.FloatValue: return PackFloat(writer, value.(*reflect.FloatValue));
    case *reflect.StringValue: return PackString(writer, value.(*reflect.StringValue));
    case *reflect.ArrayValue: return PackArray(writer, value.(*reflect.ArrayValue))
    case *reflect.SliceValue: return PackSlice(writer, value.(*reflect.SliceValue))
    case *reflect.MapValue: return PackMap(writer, value.(*reflect.MapValue))
    }
    panic("unsupported type: " + value.Type().String())
}

func Pack(writer io.Writer, value interface{}) (n int, err os.Error) {
    if value == nil { return PackNil(writer) }
	return PackValue(writer, reflect.NewValue(value));
}
