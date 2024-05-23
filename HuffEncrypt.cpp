#include <iostream>
#include <vector>
#include <string>
#include <map>
#include <algorithm>
#include <fstream>

using namespace std;

struct Node {
    char data;
    unsigned freq;
    Node *left;
    Node *right;
    bool isParent;

    Node(char _data, int _freq) : data(_data), freq(_freq), left(nullptr), right(nullptr), isParent(false) {}
};

// Criando um functor como logica para o sort. Utilizando const para evitar alteracoes indesejadas.
struct SortNodes {
    bool operator()(const Node* a, const Node* b) const{
        return a->freq > b->freq;
    }
};

Node* BuildHuffmanTree(const map<char,int>& frequencies) {
    vector<Node*> nodes;
    for(const auto& f : frequencies) {
        nodes.push_back(new Node(f.first, f.second));
    }

    sort(nodes.begin(), nodes.end(), SortNodes());

    // Cada no criado eh organizado em folhas.
    while(nodes.size() > 1) {
        Node* left = nodes.back();
        nodes.pop_back();
        Node* right = nodes.back();
        nodes.pop_back();

        Node* parent = new Node('$', left->freq + right->freq);
        parent->left = left;
        parent->right = right;
        parent->isParent = true;

        nodes.push_back(parent);
        sort(nodes.begin(), nodes.end(), SortNodes());
    }

    // Retornando a 'head' da arvore
    return nodes[0];
}

void BuildHuffmanCodes(Node* root, map<char, string>& codes, string code = "") {
    if(root == nullptr){
        return;
    }

    if(!root->isParent) {
        codes[root->data] = code;
    } else {
        BuildHuffmanCodes(root->left, codes, code + "0");
        BuildHuffmanCodes(root->right, codes, code + "1");
    }
}

void GetElements(const string& book, map<char,int>& frequencies) {
    ifstream inFile(book);
    if(!inFile) {
        cerr << "Error: could not open the book.\n";
        return;
    }

    char c;
    while(inFile.get(c)) {
        frequencies[c]++;
    }

    // Para mostrar as letras ordenadas por frequencia
    std::vector<std::pair<char, int>> sortedLetters(frequencies.begin(), frequencies.end());
    std::sort(sortedLetters.begin(), sortedLetters.end(), [](const auto& a, const auto& b) { return a.second > b.second; });

    for (const auto& entry : sortedLetters) {
        std::cout << entry.first << ": " << entry.second << std::endl;
    }
}

void EncryptFile(const string& inputFile, const string& outputFile, const map<char, string>& codes) {
    ifstream inFile(inputFile);
    if (!inFile) {
        cerr << "Error: could not open input file.\n";
        return;
    }

    ofstream outFile(outputFile);
    if(!outFile) {
        cerr << "Error: could not open output file.\n";
        inFile.close();
        return;
    }

    char c;
    while(inFile.get(c)) {
        auto it = codes.find(c);
        if(it != codes.end()) {
            outFile << it->second;
        } else {
            cerr << "Error: character nod found\n";
            inFile.close();
            outFile.close();
            return;
        }
    }

    inFile.close();
    outFile.close();
}

void freeHuffmanTree(Node* root) {
    if (root == nullptr) return;
    freeHuffmanTree(root->left);
    freeHuffmanTree(root->right);
    delete root;
}

int main() {
    string book = "book.txt";
    string input = "begin.txt";
    string output = "encrypted.txt";

    map<char, int> frequencies;
    map<char, string> codes;

    GetElements(book, frequencies);

    Node* root = BuildHuffmanTree(frequencies);
    BuildHuffmanCodes(root, codes);

    cout << "Huffman Codes:\n";
    for (const auto& entry : codes) {
        cout << "Character: " << entry.first << ", Code: " << entry.second << "\n";
    }

    EncryptFile("begin.txt", "encrypted.txt", codes);

    cout << "Encrypted file created successfully!\n";

    // Liberar a memória alocada para a árvore de Huffman
    freeHuffmanTree(root);

    return 0;
}