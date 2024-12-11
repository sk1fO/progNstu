#include<iostream>
#include<sstream>

#include "linkedlist.h"
#include "database.h"


using namespace std;


int main(){
    string db_name;
    int tuples_limit;
    table_init(db_name, tuples_limit);

    while (1){
        string input;
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
            cout << "comm Insert" << endl;

            string table = inputString[2] -> val;
            cout << table << endl;

            list values;
            int pos = input.find('(');
            string value;
            string inputValue = input.substr(pos + 1, input.length() - pos - 2);
            while ((pos = inputValue.find(", ")) != string::npos) {
                value = inputValue.substr(0, pos);
                values.push_back(value);
                inputValue.erase(0, pos + 2);
            }
            values.push_back(inputValue);


            for(int i; i < values.size(); i++){
                values[i] -> val = values[i] -> val.substr(1, values[i] -> val.length()-2);
            }
            values.print();

            string writeVal;
            for(int i; i < values.size(); i++){
                writeVal = writeVal + values[i] -> val + ",";
            }
            writeVal.pop_back();

            cout << writeVal << endl;



            if (INSERT_INTO(writeVal, table, db_name ,tuples_limit)) cout << "Sucecces" << endl;
            else cout << "Syntax error" << endl;

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