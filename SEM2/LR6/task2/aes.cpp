#include <iostream>
#include <vector>
#include <fstream>
#include <random>
#include <iomanip>
#include <cstring>

using namespace std;

// Таблица замен S-box
const unsigned char s_box[256] = {
    0x63, 0x7C, 0x77, 0x7B, 0xF2, 0x6B, 0x6F, 0xC5, 0x30, 0x01, 0x67, 0x2B, 0xFE, 0xD7, 0xAB, 0x76,
    0xCA, 0x82, 0xC9, 0x7D, 0xFA, 0x59, 0x47, 0xF0, 0xAD, 0xD4, 0xA2, 0xAF, 0x9C, 0xA4, 0x72, 0xC0,
    0xB7, 0xFD, 0x93, 0x26, 0x36, 0x3F, 0xF7, 0xCC, 0x34, 0xA5, 0xE5, 0xF1, 0x71, 0xD8, 0x31, 0x15,
    0x04, 0xC7, 0x23, 0xC3, 0x18, 0x96, 0x05, 0x9A, 0x07, 0x12, 0x80, 0xE2, 0xEB, 0x27, 0xB2, 0x75,
    0x09, 0x83, 0x2C, 0x1A, 0x1B, 0x6E, 0x5A, 0xA0, 0x52, 0x3B, 0xD6, 0xB3, 0x29, 0xE3, 0x2F, 0x84,
    0x53, 0xD1, 0x00, 0xED, 0x20, 0xFC, 0xB1, 0x5B, 0x6A, 0xCB, 0xBE, 0x39, 0x4A, 0x4C, 0x58, 0xCF,
    0xD0, 0xEF, 0xAA, 0xFB, 0x43, 0x4D, 0x33, 0x85, 0x45, 0xF9, 0x02, 0x7F, 0x50, 0x3C, 0x9F, 0xA8,
    0x51, 0xA3, 0x40, 0x8F, 0x92, 0x9D, 0x38, 0xF5, 0xBC, 0xB6, 0xDA, 0x21, 0x10, 0xFF, 0xF3, 0xD2,
    0xCD, 0x0C, 0x13, 0xEC, 0x5F, 0x97, 0x44, 0x17, 0xC4, 0xA7, 0x7E, 0x3D, 0x64, 0x5D, 0x19, 0x73,
    0x60, 0x81, 0x4F, 0xDC, 0x22, 0x2A, 0x90, 0x88, 0x46, 0xEE, 0xB8, 0x14, 0xDE, 0x5E, 0x0B, 0xDB,
    0xE0, 0x32, 0x3A, 0x0A, 0x49, 0x06, 0x24, 0x5C, 0xC2, 0xD3, 0xAC, 0x62, 0x91, 0x95, 0xE4, 0x79,
    0xE7, 0xC8, 0x37, 0x6D, 0x8D, 0xD5, 0x4E, 0xA9, 0x6C, 0x56, 0xF4, 0xEA, 0x65, 0x7A, 0xAE, 0x08,
    0xBA, 0x78, 0x25, 0x2E, 0x1C, 0xA6, 0xB4, 0xC6, 0xE8, 0xDD, 0x74, 0x1F, 0x4B, 0xBD, 0x8B, 0x8A,
    0x70, 0x3E, 0xB5, 0x66, 0x48, 0x03, 0xF6, 0x0E, 0x61, 0x35, 0x57, 0xB9, 0x86, 0xC1, 0x1D, 0x9E,
    0xE1, 0xF8, 0x98, 0x11, 0x69, 0xD9, 0x8E, 0x94, 0x9B, 0x1E, 0x87, 0xE9, 0xCE, 0x55, 0x28, 0xDF,
    0x8C, 0xA1, 0x89, 0x0D, 0xBF, 0xE6, 0x42, 0x68, 0x41, 0x99, 0x2D, 0x0F, 0xB0, 0x54, 0xBB, 0x16,
};

// Обратная таблица замен S-box
const unsigned char inv_s_box[256] = {
    0x52, 0x09, 0x6A, 0xD5, 0x30, 0x36, 0xA5, 0x38, 0xBF, 0x40, 0xA3, 0x9E, 0x81, 0xF3, 0xD7, 0xFB,
    0x7C, 0xE3, 0x39, 0x82, 0x9B, 0x2F, 0xFF, 0x87, 0x34, 0x8E, 0x43, 0x44, 0xC4, 0xDE, 0xE9, 0xCB,
    0x54, 0x7B, 0x94, 0x32, 0xA6, 0xC2, 0x23, 0x3D, 0xEE, 0x4C, 0x95, 0x0B, 0x42, 0xFA, 0xC3, 0x4E,
    0x08, 0x2E, 0xA1, 0x66, 0x28, 0xD9, 0x24, 0xB2, 0x76, 0x5B, 0xA2, 0x49, 0x6D, 0x8B, 0xD1, 0x25,
    0x72, 0xF8, 0xF6, 0x64, 0x86, 0x68, 0x98, 0x16, 0xD4, 0xA4, 0x5C, 0xCC, 0x5D, 0x65, 0xB6, 0x92,
    0x6C, 0x70, 0x48, 0x50, 0xFD, 0xED, 0xB9, 0xDA, 0x5E, 0x15, 0x46, 0x57, 0xA7, 0x8D, 0x9D, 0x84,
    0x90, 0xD8, 0xAB, 0x00, 0x8C, 0xBC, 0xD3, 0x0A, 0xF7, 0xE4, 0x58, 0x05, 0xB8, 0xB3, 0x45, 0x06,
    0xD0, 0x2C, 0x1E, 0x8F, 0xCA, 0x3F, 0x0F, 0x02, 0xC1, 0xAF, 0xBD, 0x03, 0x01, 0x13, 0x8A, 0x6B,
    0x3A, 0x91, 0x11, 0x41, 0x4F, 0x67, 0xDC, 0xEA, 0x97, 0xF2, 0xCF, 0xCE, 0xF0, 0xB4, 0xE6, 0x73,
    0x96, 0xAC, 0x74, 0x22, 0xE7, 0xAD, 0x35, 0x85, 0xE2, 0xF9, 0x37, 0xE8, 0x1C, 0x75, 0xDF, 0x6E,
    0x47, 0xF1, 0x1A, 0x71, 0x1D, 0x29, 0xC5, 0x89, 0x6F, 0xB7, 0x62, 0x0E, 0xAA, 0x18, 0xBE, 0x1B,
    0xFC, 0x56, 0x3E, 0x4B, 0xC6, 0xD2, 0x79, 0x20, 0x9A, 0xDB, 0xC0, 0xFE, 0x78, 0xCD, 0x5A, 0xF4,
    0x1F, 0xDD, 0xA8, 0x33, 0x88, 0x07, 0xC7, 0x31, 0xB1, 0x12, 0x10, 0x59, 0x27, 0x80, 0xEC, 0x5F,
    0x60, 0x51, 0x7F, 0xA9, 0x19, 0xB5, 0x4A, 0x0D, 0x2D, 0xE5, 0x7A, 0x9F, 0x93, 0xC9, 0x9C, 0xEF,
    0xA0, 0xE0, 0x3B, 0x4D, 0xAE, 0x2A, 0xF5, 0xB0, 0xC8, 0xEB, 0xBB, 0x3C, 0x83, 0x53, 0x99, 0x61,
    0x17, 0x2B, 0x04, 0x7E, 0xBA, 0x77, 0xD6, 0x26, 0xE1, 0x69, 0x14, 0x63, 0x55, 0x21, 0x0C, 0x7D,
};

// Константы Rcon для расширения ключа
const unsigned char rcon[11] = {
    0x00, 0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80, 0x1B, 0x36
};

// Расширение ключа: генерирует раундовые ключи из начального ключа
void expand_key(const unsigned char* key, unsigned char* roundKeys) {
    // Копируем начальный ключ в первые 16 байтов roundKeys
    memcpy(roundKeys, key, 16);
    unsigned char temp[4]; // Временный массив для хранения 4 байт текущего слова
    int wordIndex = 4; // Счётчик слов (4 байта на слово, начальное значение 4, так как первые 4 слова это начальный ключ)

    // Цикл для генерации всех раундовых ключей
    while (wordIndex < 4 * (10 + 1)) { // Количество раундовых ключей зависит от числа раундов 
        // Копируем последние 4 байта предыдущего раундового ключа в temp
        memcpy(temp, roundKeys + (wordIndex - 1) * 4, 4);
        if (wordIndex % 4 == 0) { // Каждое четвёртое слово обрабатывается отдельно
            // Циклический сдвиг слов
            unsigned char t = temp[0];
            temp[0] = s_box[temp[1]] ^ rcon[wordIndex / 4]; // Применение S-box и Rcon
            temp[1] = s_box[temp[2]];
            temp[2] = s_box[temp[3]];
            temp[3] = s_box[t];
        }
        // XOR текущего слова с соответствующим словом из предыдущих 16 байт
        for (int j = 0; j < 4; ++j) {
            roundKeys[wordIndex * 4 + j] = roundKeys[(wordIndex - 4) * 4 + j] ^ temp[j];
        }
        ++wordIndex; // Увеличиваем счётчик слов
    }
}

// Добавление раундового ключа к состоянию
void add_round_key(unsigned char state[4][4], const unsigned char* roundKey) {
    for (int row = 0; row < 4; ++row) { // Проходим по строкам состояния
        for (int col = 0; col < 4; ++col) { // Проходим по столбцам состояния
            state[row][col] ^= roundKey[row + 4 * col]; // Добавляем (XOR) байт раундового ключа к байту состояния
        }
    }
}

// Замена байтов в состоянии с использованием S-box
void sub_bytes(unsigned char state[4][4]) {
    for (int row = 0; row < 4; ++row) { // Проходим по строкам состояния
        for (int col = 0; col < 4; ++col) { // Проходим по столбцам состояния
            state[row][col] = s_box[state[row][col]]; // Заменяем байт с использованием таблицы S-box
                    }
    }
}

// Обратная замена байтов в состоянии с использованием обратной таблицы S-box
void inv_sub_bytes(unsigned char state[4][4]) {
    for (int row = 0; row < 4; ++row) { // Проходим по строкам состояния
        for (int col = 0; col < 4; ++col) { // Проходим по столбцам состояния
            state[row][col] = inv_s_box[state[row][col]]; // Заменяем байт с использованием обратной таблицы S-box
        }
    }
}

// Сдвиг строк состояния
void shift_rows(unsigned char state[4][4]) {
    unsigned char temp; // Временная переменная для хранения байта

    // Сдвиг второй строки на 1 байт влево
    temp = state[1][0];
    state[1][0] = state[1][1];
    state[1][1] = state[1][2];
    state[1][2] = state[1][3];
    state[1][3] = temp;

    // Сдвиг третьей строки на 2 байта влево
    temp = state[2][0];
    state[2][0] = state[2][2];
    state[2][2] = temp;
    temp = state[2][1];
    state[2][1] = state[2][3];
    state[2][3] = temp;

    // Сдвиг четвертой строки на 3 байта влево
    temp = state[3][0];
    state[3][0] = state[3][3];
    state[3][3] = state[3][2];
    state[3][2] = state[3][1];
    state[3][1] = temp;
}

// Обратный сдвиг строк состояния
void inv_shift_rows(unsigned char state[4][4]) {
    unsigned char temp; // Временная переменная для хранения байта

    // Сдвиг второй строки на 1 байт вправо
    temp = state[1][3];
    state[1][3] = state[1][2];
    state[1][2] = state[1][1];
    state[1][1] = state[1][0];
    state[1][0] = temp;

    // Сдвиг третьей строки на 2 байта вправо
    temp = state[2][3];
    state[2][3] = state[2][1];
    state[2][1] = temp;
    temp = state[2][2];
    state[2][2] = state[2][0];
    state[2][0] = temp;

    // Сдвиг четвертой строки на 3 байта вправо
    temp = state[3][3];
    state[3][3] = state[3][0];
    state[3][0] = state[3][1];
    state[3][1] = state[3][2];
    state[3][2] = temp;
}

// Умножение в поле Галуа
unsigned char xtime(unsigned char x) {
    return ((x << 1) & 0xFE) ^ (((x >> 7) & 1) * 0x1B); // Умножение на x (0x02) в поле Галуа GF(2^8)
}

// Перемешивание столбцов состояния
void mix_columns(unsigned char state[4][4]) {
    for (int col = 0; col < 4; ++col) { // Проходим по каждому столбцу состояния
        unsigned char originalBytes[4]; // Сохраняем оригинальные байты столбца
        unsigned char multipliedBytes[4]; // Сохраняем результат умножения байтов на 0x02 (x)
        for (int row = 0; row < 4; ++row) {
            originalBytes[row] = state[row][col];
            multipliedBytes[row] = xtime(originalBytes[row]);
        }
        // Применяем умножение и сложение в поле Галуа для каждого байта столбца
        state[0][col] = multipliedBytes[0] ^ originalBytes[1] ^ multipliedBytes[1] ^ originalBytes[2] ^ originalBytes[3];
        state[1][col] = originalBytes[0] ^ multipliedBytes[1] ^ originalBytes[2] ^ multipliedBytes[2] ^ originalBytes[3];
        state[2][col] = originalBytes[0] ^ originalBytes[1] ^ multipliedBytes[2] ^ originalBytes[3] ^ multipliedBytes[3];
        state[3][col] = originalBytes[0] ^ multipliedBytes[0] ^ originalBytes[1] ^ originalBytes[2] ^ multipliedBytes[3];
    }
}

// Обратное перемешивание столбцов состояния
void inv_mix_columns(unsigned char state[4][4]) {
    for (int col = 0; col < 4; ++col) { // Проходим по каждому столбцу состояния
        unsigned char originalBytes[4]; // Сохраняем оригинальные байты столбца
        unsigned char multipliedBytes[4]; // Сохраняем результат умножения байтов на 0x02 (x)
        for (int row = 0; row < 4; ++row) {
            originalBytes[row] = state[row][col];
            multipliedBytes[row] = xtime(xtime(originalBytes[row] ^ originalBytes[(row + 1) % 4]));
        }
        // Применяем умножение и сложение в поле Галуа для каждого байта столбца
        state[0][col] = multipliedBytes[0] ^ multipliedBytes[1] ^ originalBytes[2] ^ originalBytes[3];
        state[1][col] = originalBytes[0] ^ multipliedBytes[1] ^ multipliedBytes[2] ^ originalBytes[3];
        state[2][col] = originalBytes[0] ^ originalBytes[1] ^ multipliedBytes[2] ^ multipliedBytes[3];
        state[3][col] = multipliedBytes[0] ^ originalBytes[1] ^ originalBytes[2] ^ multipliedBytes[3];
    }
}

// Основная функция шифрования блока
void encrypt_block(unsigned char* input, unsigned char* output, const unsigned char* roundKeys) {
    unsigned char state[4][4]; // Матрица состояния
    // Копирование входного блока в состояние
    for (int row = 0; row < 4; ++row) {
        for (int col = 0; col < 4; ++col) {
            state[row][col] = input[row + 4 * col]; // Формируем матрицу состояния из входных данных
        }
    }

    // Начальное добавление раундового ключа
    add_round_key(state, roundKeys);

    // Основные раунды
    for (int round = 1; round < 10; ++round) {
        sub_bytes(state); // Замена байтов
        shift_rows(state); // Сдвиг строк
        mix_columns(state); // Перемешивание столбцов
        add_round_key(state, roundKeys + round * 16); // Добавление раундового ключа
    }

    // Последний раунд без перемешивания столбцов
    sub_bytes(state);
    shift_rows(state);
    add_round_key(state, roundKeys + 10 * 16);

    // Копирование состояния в выходной блок
    for (int row = 0; row < 4; ++row) {
        for (int col = 0; col < 4; ++col) {
            output[row + 4 * col] = state[row][col]; // Переносим зашифрованные данные в выходной массив
        }
    }
}

// Основная функция дешифрования блока
void decrypt_block(unsigned char* input, unsigned char* output, const unsigned char* roundKeys) {
    unsigned char state[4][4]; // Матрица состояния
    // Копирование входного блока в состояние
    for (int row = 0; row < 4; ++row) {
        for (int col = 0; col < 4; ++col) {
            state[row][col] = input[row + 4 * col]; // Формируем матрицу состояния из входных данных
        }
    }

    // Начальное добавление раундового ключа
    add_round_key(state, roundKeys + 10 * 16);

    // Основные раунды
    for (int round = 10 - 1; round > 0; --round) {
        inv_shift_rows(state); // Обратный сдвиг строк
        inv_sub_bytes(state); // Обратная замена байтов
        add_round_key(state, roundKeys + round * 16); // Добавление раундового ключа
        inv_mix_columns(state); // Обратное перемешивание столбцов
    }

    // Последний раунд без обратного перемешивания столбцов
    inv_shift_rows(state);
    inv_sub_bytes(state);
    add_round_key(state, roundKeys);

    // Копирование состояния в выходной блок
    for (int row = 0; row < 4; ++row) {
        for (int col = 0; col < 4; ++col) {
            output[row + 4 * col] = state[row][col]; // Переносим расшифрованные данные в выходной массив
        }
    }
}

// Функция шифрования в режиме CFB
void encrypt_cfb(const vector<unsigned char>& plaintext, const vector<unsigned char>& key, const vector<unsigned char>& iv, vector<unsigned char>& ciphertext) {
    unsigned char roundKeys[176]; // Массив для хранения всех раундовых ключей
    expand_key(key.data(), roundKeys); // Генерация раундовых ключей из начального ключа

    unsigned char iv_copy[16]; // Копия вектора инициализации для модификации в процессе шифрования
    memcpy(iv_copy, iv.data(), 16); // Копируем начальный вектор инициализации

    ciphertext.resize(plaintext.size()); // Изменяем размер вектора ciphertext, чтобы он соответствовал размеру plaintext
    for (size_t i = 0; i < plaintext.size(); i += 16) { // Проходим по plaintext блоками по 16 байт
        unsigned char outputBlock[16]; // Временный массив для хранения зашифрованного блока
        encrypt_block(iv_copy, outputBlock, roundKeys); // Шифруем текущий вектор инициализации

        for (size_t j = 0; j < 16 && (i + j) < plaintext.size(); ++j) { // Проходим по текущему блоку
            ciphertext[i + j] = plaintext[i + j] ^ outputBlock[j]; // XOR с plaintext для получения ciphertext
            iv_copy[j] = ciphertext[i + j]; // Обновляем вектор инициализации
        }
    }
}

// Функция дешифрования в режиме CFB
void decrypt_cfb(const vector<unsigned char>& ciphertext, const vector<unsigned char>& key, const vector<unsigned char>& iv, vector<unsigned char>& plaintext) {
    unsigned char roundKeys[176]; // Массив для хранения всех раундовых ключей
    expand_key(key.data(), roundKeys); // Генерация раундовых ключей из начального ключа

    unsigned char iv_copy[16]; // Копия вектора инициализации для модификации в процессе дешифрования
    memcpy(iv_copy, iv.data(), 16); // Копируем начальный вектор инициализации

    plaintext.resize(ciphertext.size()); // Изменяем размер вектора plaintext, чтобы он соответствовал размеру ciphertext
    for (size_t i = 0; i < ciphertext.size(); i += 16) { // Проходим по ciphertext блоками по 16 байт
        unsigned char outputBlock[16]; // Временный массив для хранения зашифрованного блока
        encrypt_block(iv_copy, outputBlock, roundKeys); // Шифруем текущий вектор инициализации

        for (size_t j = 0; j < 16 && (i + j) < ciphertext.size(); ++j) { // Проходим по текущему блоку
            plaintext[i + j] = ciphertext[i + j] ^ outputBlock[j]; // XOR с ciphertext для получения plaintext
            iv_copy[j] = ciphertext[i + j]; // Обновляем вектор инициализации
        }
    }
}

// Функция для вывода раундовых ключей
void print_round_keys(const unsigned char* roundKeys) {
    cout << endl << "Раундовые ключи:" << endl;
    for (int round = 0; round <= 10; ++round) {
        cout << "Раунд " << round << ":" << endl;
        for (int i = 0; i < 16; ++i) {
            cout << hex << setw(2) << setfill('0') << static_cast<int>(roundKeys[round * 16 + i]) << " "; // Вывод байта в шестнадцатеричном формате
            if ((i + 1) % 4 == 0) cout << endl; // Переход на новую строку после каждых 4 байт
        }
        cout << endl;
    }
}

// Функция для записи данных в файл
void write_output_to_file(const string& filename, const unsigned char* roundKeys, const vector<unsigned char>& ciphertext, const string& decrypted_str) {
    ofstream output_file(filename); // Открываем файл для записи
    if (!output_file) { // Проверка на ошибку открытия файла
        cerr << "Ошибка открытия файла для записи." << endl;
        return;
    }

    // Запись раундовых ключей
    output_file << "Раундовые ключи:" << endl;
    for (int round = 0; round <= 10; ++round) {
        output_file << "Раунд " << round << ":" << endl;
        for (int i = 0; i < 16; ++i) {
            output_file << hex << setw(2) << setfill('0') << static_cast<int>(roundKeys[round * 16 + i]) << " "; // Запись байта в шестнадцатеричном формате
            if ((i + 1) % 4 == 0) output_file << endl; // Переход на новую строку после каждых 4 байт
        }
        output_file << endl;
    }

    // Запись зашифрованного текста
    output_file << endl << "Зашифрованный текст: ";
    for (auto c : ciphertext) {
        output_file << hex << static_cast<int>(c); // Запись зашифрованного текста в шестнадцатеричном формате
    }
    output_file << endl;

    // Запись расшифрованного текста
    output_file << "Расшифрованный текст: " << decrypted_str << endl; // Запись расшифрованного текста

    output_file.close(); // Закрываем файл
}
int main() {
    vector<unsigned char> key; // Вектор для хранения ключа
    vector<unsigned char> iv; // Вектор для хранения вектора инициализации

    string key_str; // Строка для хранения мастер-ключа
    cout << "Введите мастер-ключ (16 символов): ";
    cin >> key_str;
    if (key_str.size() != 16) { // Проверка длины мастер-ключа
        cerr << "Ошибка: Мастер-ключ должен состоять из 16 символов." << endl;
        return 1;
    }
    key.assign(key_str.begin(), key_str.end()); // Преобразуем строку в вектор байтов

    // Инициализация генератора случайных чисел непосредственно перед использованием
    mt19937_64 mt(random_device{}()); // Инициализация генератора случайных чисел
    uniform_int_distribution<int> dis(0, 255); // Диапазон для генерации случайных чисел
    for (int i = 0; i < 16; ++i) {
        iv.push_back(dis(mt)); // Генерация вектора инициализации
    }

    string plaintext; // Строка для хранения текста для шифрования
    cout << "Введите текст для шифрования: ";
    cin.ignore();
    getline(cin, plaintext);
    vector<unsigned char> plaintext_vec(plaintext.begin(), plaintext.end()); // Преобразуем строку в вектор байтов

    vector<unsigned char> ciphertext; // Вектор для хранения зашифрованного текста
    unsigned char roundKeys[176]; // Массив для хранения всех раундовых ключей
    expand_key(key.data(), roundKeys); // Генерация раундовых ключей из начального ключа

    print_round_keys(roundKeys); // Вывод раундовых ключей

    encrypt_cfb(plaintext_vec, key, iv, ciphertext); // Шифрование текста в режиме CFB
    cout << endl << "Зашифрованный текст: ";
    for (auto c : ciphertext) {
        cout << hex << static_cast<int>(c); // Вывод зашифрованного текста в шестнадцатеричном формате
    }
    cout << endl;

    vector<unsigned char> decryptedtext; // Вектор для хранения расшифрованного текста
    decrypt_cfb(ciphertext, key, iv, decryptedtext); // Дешифрование текста в режиме CFB

    string decrypted_str(decryptedtext.begin(), decryptedtext.end()); // Преобразуем вектор байтов в строку
    cout << endl << "Расшифрованный текст: " << decrypted_str << endl; // Вывод расшифрованного текста

    cout << endl << "Данные записаны в output.txt" << endl;
    write_output_to_file("output.txt", roundKeys, ciphertext, decrypted_str); // Запись данных в файл

    return 0; 
}