//#include <string>
#include "linkedlist.h"
#include <sstream>
#include <iostream>

using namespace std;



int main()
{

    while (1){
        string input = "";
        cout << "Type command: \n> ";
        getline(cin, input);
        list inputString;
        istringstream is(input);
        string tmp;
        
        while (is >> tmp)
        {
            inputString.push_back(tmp);
        } 
        
        

        
        if(inputString[0]->val == "DELETE"){
            cout << "comm Delete" << endl;
            
            string tableColl = inputString[4]->val;
            int pos = tableColl.find('.');
            string table = tableColl.substr(0, pos);
            string coll = tableColl.substr(pos + 1, tableColl.length());
            
            cout << coll << table << endl;

            
            pos = input.find("'");

            string value = input.substr(pos + 1, input.length() - pos - 2);


            cout << value << endl;

        }


        else if (inputString[0] -> val == "INSERT"){
            inputString.print();
            cout << "comm Insert" << endl;

            cout << input << endl;

            string table = inputString[2] -> val;
            cout << table << endl;

            list values;
            int pos = input.find('(');
            string value;
            string inputValue = input.substr(pos + 1, input.length() - pos - 2);

            cout << inputValue << endl;
            cout << pos << endl;
            pos = 0;
            while ((pos = inputValue.find(", ")) != string::npos) {
                cout << inputValue << endl;
                value = inputValue.substr(1, inputValue.length() - pos -1);
                values.push_back(value);
                inputValue.erase(0, pos + 2);
            }
            values.push_back(inputValue);


            for(int i; i < values.size(); i++){
                values[i] -> val = values[i] -> val.substr(1, values[i] -> val.length()-2);
            }
            values.print();

            for (int i; i < values.size(); i++) values.remove_last();

        }

        else if (inputString[0]-> val == "SELECT"){
            cout << "comm Select" << endl;


            if (input.find("WHERE") != string::npos){
                cout << "where token" << endl;


            }
        }





        else if (inputString[0]-> val == "EXIT"){
            cout << "EXIT" << endl;
            break;
        }
        
        else {
            cout << "Wrong command" << endl;
        }

        

    }
    
}