#include <iostream>
#include <string>
#include <vector>


using namespace std;

class Device { // абстрактный класс устройств
public:
	// виртуальные функции объявлены равными нулю, что делает класс абстрактным (это означает, что класс нельзя создать напрямую, а только наследовать от него)
	virtual void getName() const = 0;
	virtual void getInfo() const = 0;
	virtual string getNumber() const = 0;
protected:
	string brand;
	int price;
	string nomenclatureNumber;
};

// конкретные классы устройств
class Headphones : public Device {
private:
	string constructionType;
	string attachmentMethod;

public:
	// реализация виртуальных функций для наушников
	void getName() const override {
		cout << nomenclatureNumber << " " << brand << " " << price << endl;
	}

	void getInfo() const override {
		cout << brand << " " << price << " " << nomenclatureNumber << " " << constructionType << " " << attachmentMethod << endl;
	}

	string getNumber() const override {
		return nomenclatureNumber;
	}

	Headphones(string brand, int price, string nomenclatureNumber, string constructionType, string attachmentMethod) {
		this->brand = brand;
		this->price = price;
		this->nomenclatureNumber = nomenclatureNumber;
		this->constructionType = constructionType;
		this->attachmentMethod = attachmentMethod;
	}
};

class Microphones : public Device {
private:
	string frequencyRange;
	string sensitivity;

public:
	// реализация виртуальных функций для микрофонов
	void getName() const override {
		cout << nomenclatureNumber << " " << brand << " " << price << endl;
	}

	void getInfo() const override {
		cout << brand << " " << price << " " << nomenclatureNumber << " " << frequencyRange << "Hz " << sensitivity << "dB/w/m" << endl;
	}

	string getNumber() const override {
		return nomenclatureNumber;
	}

	Microphones(string brand, int price, string nomenclatureNumber, string frequencyRange, string sensitivity) {
		this->brand = brand;
		this->price = price;
		this->nomenclatureNumber = nomenclatureNumber;
		this->frequencyRange = frequencyRange;
		this->sensitivity = sensitivity;
	}
};

class Keyboards : public Device {
private:
	string switchType;
	string interfaceKeybord;

public:
	// реализация виртуальных функций для клавиатур
	void getName() const override {
		cout << nomenclatureNumber << " " << brand << " " << price << endl;
	}

	void getInfo() const override {
		cout << brand << " " << price << " " << nomenclatureNumber << " " << switchType << " " << interfaceKeybord << endl;
	}

	string getNumber() const override {
		return nomenclatureNumber;
	}

	Keyboards(string brand, int price, string nomenclatureNumber, string switchType, string interfaceKeybord) {
		this->brand = brand;
		this->price = price;
		this->nomenclatureNumber = nomenclatureNumber;
		this->switchType = switchType;
		this->interfaceKeybord = interfaceKeybord;
	}
};

// интерфейс "фабрики" устройств
class DeviceFactory {
public:
	virtual Device* createDevice(string, int, string, string, string) = 0; // виртуальная функция createDevice() объявлена равной нулю, что делает класс абстрактным
};

// конкретные классы "фабрики" устройств
class HeadphonesFactory : public DeviceFactory {
public:
	// реализация виртуальной функции createDevice() для создания наушников
	Device* createDevice(string brand, int price, string nomenclatureNumber, string constructionType, string attachmentMethod) override {
		return new Headphones(brand, price, nomenclatureNumber, constructionType, attachmentMethod);
	}
};

class MicrophonesFactory : public DeviceFactory {
public:
	Device* createDevice(string brand, int price, string nomenclatureNumber, string frequencyRange, string sensitivity) override {
		return new Microphones(brand, price, nomenclatureNumber, frequencyRange, sensitivity);
	}
};
class KeyboardsFactory : public DeviceFactory {
public:
	Device* createDevice(string brand, int price, string nomenclatureNumber, string switchType, string interfaceKeybord) override {
		return new Keyboards(brand, price, nomenclatureNumber, switchType, interfaceKeybord);
	}
};

int main() {
	// использование фабричного метода для создания устройств
	vector<Device*> devices;

	// создание фабрик для наушников
	DeviceFactory* headphonesFactory = new HeadphonesFactory();
    DeviceFactory* microphonesFactory = new MicrophonesFactory();
    DeviceFactory* keyboardsFactory = new KeyboardsFactory();

	// создание наушников с использованием фабрики и вызов его методов
	Device* sonyHeadphones = headphonesFactory->createDevice("Sony", 1500, "h1", "Вкладыши", "Дуговое крепление");
    devices.push_back(sonyHeadphones);

    Device* marchallHeadphones = headphonesFactory->createDevice("Marchall", 3000, "h2", "Накладные", "Дуговое крепление");
    devices.push_back(marchallHeadphones);

	
	Device* MicroSoundMic = microphonesFactory->createDevice("MicroSound", 980, "m1", "140-5000", "100");
	devices.push_back(MicroSoundMic);

    Device* HiperMic = microphonesFactory->createDevice("Hiper", 5000, "m2", "20-20000", "105");
	devices.push_back(HiperMic);

	
	Device* logitechKeyboard = keyboardsFactory->createDevice("Logitech", 1250, "k1", "Линейные переключатели", "USB");
	devices.push_back(logitechKeyboard);

    Device* razerKeyboard = keyboardsFactory->createDevice("Razer", 3200, "k2", "Линейные переключатели", "USB + Bluetooth");
	devices.push_back(razerKeyboard);

    
    system("clear");
	int choice;
	do {
		cout << "Выберите действие из предложенных:" << endl;
		cout << "1. Вывести полный список устройств" << endl;
		cout << "2. Вывести информацию по устройству" << endl;
		cout << "3. Выйти из программы" << endl << ">>> ";
		cin >> choice;
		string nomenclatureNumber;
		switch (choice) {
		case 1:
            system("clear");
			for (auto item : devices) {
				item->getName();
			}
            cout << endl;
			break;

		case 2:
			cout << "Введите номенклатурный номер устройства: ";
			cin >> nomenclatureNumber;
			for (auto item : devices) {
				if (item->getNumber() == nomenclatureNumber) {
					item->getInfo();
					break;
				}
			}
			break;
		case 3:
			cout << "Выход" << endl;
			break;
		default:
			cout << "Некорректный выбор" << endl;
		}
	} while (choice != 3);

	return 0;
}
