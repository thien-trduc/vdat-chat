package utils

//var ArraySocketId []string
//
//func ReadLines(path string) ([]string, error) {
//	file, err := os.Open(path)
//	if err != nil {
//		return nil, err
//	}
//	defer file.Close()
//
//	var lines []string
//	scanner := bufio.NewScanner(file)
//	for scanner.Scan() {
//		lines = append(lines, scanner.Text())
//	}
//	return lines, scanner.Err()
//}
//
//// writeLines writes the lines to the given file.
//func WriteLines(lines []string, path string) error {
//	file, err := os.Create(path)
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//
//	w := bufio.NewWriter(file)
//	for _, line := range lines {
//		fmt.Fprintln(w, line)
//	}
//	return w.Flush()
//}
//func CreateArrayId() []string {
//	var i int
//	for i = 1; i < 10000; i++ {
//		ArraySocketId = append(ArraySocketId, strconv.Itoa(i))
//
//	}
//	return ArraySocketId
//}
//func CheckFileSocketId() error {
//	ArraySocketId, _ = ReadLines("socketid.data")
//
//	if len(ArraySocketId) == 0 {
//		CreateArrayId()
//	}
//	return nil
//}
//func DeleteItemInArray(a []string) []string {
//	a[0] = a[len(a)-1] // Copy last element to index i.
//	a[len(a)-1] = ""   // Erase last element (write zero value).
//	a = a[:len(a)-1]   // Truncate slice.
//	return a
//}
//func RestoreItemArray(a []string, item string) []string {
//	a = append(a, item)
//	return a
//}
//func Test() {
//	fmt.Println("end")
//}
