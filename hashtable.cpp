#include <iostream>

using namespace std;

template<typename K, typename V>
class HashTable {
private:
    struct Node {
        K key;
        V value;
        bool active; // Для отслеживания удаленных элементов

        Node() : active(false) {}
        Node(const K& k, const V& v) : key(k), value(v), active(true) {}
    };

    Node* table;
    int size;
    int capacity;

    int HashFunction(const K& key) const {
        return hash<K>()(key) % capacity;
    }

public:
    HashTable(int cap = 8) : capacity(cap), size(0) {
        table = new Node[capacity];
    }

    ~HashTable() {
        delete[] table;
    }

    void Add(const K& key, const V& value) {
        if (size >= capacity / 2) {
            Resize(); // Увеличиваем размер, если загруженность превышает 50%
        }

        int index = HashFunction(key);

        while (table[index].active) {
            if (table[index].key == key) { // Если элемент с таким ключом уже есть, обновляем значение
                table[index].value = value;
                return;
            }
            index = (index + 1) % capacity; // Линейное пробирование
        }

        table[index] = Node(key, value);
        size++;
    }

    bool Get(const K& key, V& value) const {
        int index = HashFunction(key);

        while (table[index].active) {
            if (table[index].key == key) {
                value = table[index].value;
                return true; // Значение найдено
            }
            index = (index + 1) % capacity; // Линейное пробирование
        }

        return false; // Значение не найдено
    }

    bool Remove(const K& key) {
        int index = HashFunction(key);

        while (table[index].active) {
            if (table[index].key == key) {
                table[index].active = false; // Удаляем элемент
                size--;
                return true; // Успешно удалено
            }
            index = (index + 1) % capacity; // Линейное пробирование
        }

        return false; // Элемент не найден
    }

    void Resize() {
        int old_capacity = capacity;
        capacity *= 2;
        Node* old_table = table;
        table = new Node[capacity];
        size = 0;

        for (int i = 0; i < old_capacity; ++i) {
            if (old_table[i].active) {
                Add(old_table[i].key, old_table[i].value);
            }
        }

        delete[] old_table; // Освобождаем старую таблицу
    }

    void Print() const {
        for (int i = 0; i < capacity; ++i) {
            if (table[i].active) {
                cout << "Key: " << table[i].key << ", Value: " << table[i].value << endl;
            }
        }
    }
};

int main() {
    HashTable<string, string> hash_table;

    hash_table.Add("apple", "one");
    hash_table.Add("banana", "two");
    hash_table.Add("orange", "three");

    string value;
    if (hash_table.Get("banana", value)) {
        cout << "Value for 'banana': " << value << endl;
    } else {
        cout << "'banana' not found!" << endl;
    }

    //hash_table.Remove("apple");
    hash_table.Print(); // Печатаем оставшиеся элементы

    return 0;
}