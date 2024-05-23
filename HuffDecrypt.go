package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "sort"
)

// Node represents a node in the Huffman tree.
type Node struct {
    Char     rune
    Freq     int
    Left     *Node
    Right    *Node
    IsParent bool
}

// Nodes implements sort.Interface for sorting nodes by frequency.
type Nodes []*Node

func (n Nodes) Len() int           { return len(n) }
func (n Nodes) Less(i, j int) bool { return n[i].Freq < n[j].Freq }
func (n Nodes) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }

// BuildHuffmanTree builds the Huffman tree based on character frequencies.
func BuildHuffmanTree(frequencies map[rune]int) *Node {
    var nodes Nodes
    for char, freq := range frequencies {
        nodes = append(nodes, &Node{Char: char, Freq: freq})
    }
    sort.Sort(nodes)

    for len(nodes) > 1 {
        // Combine the two nodes with the lowest frequency.
        left := nodes[0]
        right := nodes[1]
        parent := &Node{
            Freq:     left.Freq + right.Freq,
            Left:     left,
            Right:    right,
            IsParent: true,
        }
        nodes = append(nodes[2:], parent)
        sort.Sort(nodes)
    }
    return nodes[0] // Returns the root of the Huffman tree.
}

// BuildHuffmanCodes builds the Huffman codes from the Huffman tree.
func BuildHuffmanCodes(root *Node, code string, codes map[rune]string) {
    if root == nil {
        return
    }

    if !root.IsParent {
        codes[root.Char] = code
    } else {
        BuildHuffmanCodes(root.Left, code+"0", codes)
        BuildHuffmanCodes(root.Right, code+"1", codes)
    }
}

// PrintHuffmanTree prints the Huffman tree (for debugging purposes).
func PrintHuffmanTree(root *Node, code string) {
    if root == nil {
        return
    }

    if !root.IsParent {
        fmt.Printf("Character: %c, Code: %s\n", root.Char, code)
    } else {
        PrintHuffmanTree(root.Left, code+"0")
        PrintHuffmanTree(root.Right, code+"1")
    }
}

// DecryptFile decrypts a text file using the Huffman tree and writes the decrypted content to a new file.
func DecryptFile(inputFile, outputFile string, root *Node) error {
    // Read the content of the encrypted text file.
    content, err := ioutil.ReadFile(inputFile)
    if err != nil {
        return fmt.Errorf("error reading the encrypted file: %v", err)
    }

    // Create the decrypted text content.
    decryptedContent := ""
    currentNode := root
    for _, bit := range string(content) {
        if bit == '0' {
            currentNode = currentNode.Left
        } else {
            currentNode = currentNode.Right
        }

        if currentNode.Left == nil && currentNode.Right == nil {
            decryptedContent += string(currentNode.Char)
            currentNode = root
        }
    }

    // Write the decrypted content to the new file.
    err = ioutil.WriteFile(outputFile, []byte(decryptedContent), 0644)
    if err != nil {
        return fmt.Errorf("error writing the decrypted file: %v", err)
    }

    return nil
}

func main() {
    // Read the content of the book.
    content, err := ioutil.ReadFile("/home/alfredo/IME/DSA/trabHuff/book.txt")
    if err != nil {
        log.Fatalf("Error reading the file: %v", err)
    }

    // Count the frequency of each character.
    frequencies := make(map[rune]int)
    for _, char := range string(content) {
        frequencies[char]++
    }

    // Sort characters by frequency.
    sortedLetters := make([]rune, 0, len(frequencies))
    for letter := range frequencies {
        sortedLetters = append(sortedLetters, letter)
    }
    sort.Slice(sortedLetters, func(i, j int) bool {
        return frequencies[sortedLetters[i]] > frequencies[sortedLetters[j]]
    })

    // Display frequencies.
    for _, letter := range sortedLetters {
        fmt.Printf("%c: %d\n", letter, frequencies[letter])
    }

    // Build the Huffman tree.
    root := BuildHuffmanTree(frequencies)

    // Build Huffman codes.
    codes := make(map[rune]string)
    BuildHuffmanCodes(root, "", codes)

    fmt.Println("Huffman Codes:")
    for char, code := range codes {
        fmt.Printf("Character: %c, Code: %s\n", char, code)
    }

    // Decrypt the file.
    err = DecryptFile("encrypted.txt", "decrypted.txt", root)
    if err != nil {
        log.Fatalf("Error decrypting the file: %v", err)
    }
    fmt.Println("Decrypted file created successfully!")
}
