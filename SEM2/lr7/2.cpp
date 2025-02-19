#include <iostream>
#include <vector>

using namespace std;

//Если длина массива четная, то первый игрок может выиграть, взяв все четные или все нечетные числа
//Если длина массива нечетная, то первый игрок выиграет, если макс. разница между суммами четных и нечетных чисел положительна

void print(vector<int> const& nums) { //вывод вектора
    for (auto const& num : nums) {
        cout << num << ' ';
    }
    cout << endl;
}

int maximum(vector<int>& nums, int i, int j) { //вычисление максимальной разницы между суммами чисел на четных и нечетных позициях
    if (i == j) {
        return nums[i];
    }
    return max(nums[i] - maximum(nums, i + 1, j), nums[j] - maximum(nums, i, j - 1));
}

bool canWin(int n, vector<int>& nums) {
    if (n % 2 == 0) {
        return true;
    }
    int ans = maximum(nums, 0, n - 1);
    return ans >= 0;
}

void createVector() { //создание массива чисел
    int n, num;
    cout << "Введите размер массива: ";
    cin >> n;
    while (n < 1 || n > 20) {
        cout << "Размер массива недопустим. Введите число, которое >= 1, но <= 20: ";
        cin >> n;
    }

    vector <int> nums;
    cout << "Введите элементы массива: " << endl;
    while (n > 0) {
        cin >> num;

        while (num < 0 || num > 10000000) {
            cout << "Элемент массива недопустим. Введите число, которое >= 0, но <= 10 в 7 степени: ";
            cin >> num;
        }

        nums.push_back(num);
        n--;
    }

    //print(nums);

    if (canWin(nums.size(), nums) == 1) {
        cout << "true";
    }
    else cout << "false";
}

int main() {
    createVector();
    return 0;
}
