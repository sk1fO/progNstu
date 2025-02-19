#pragma once
#include <vector>
#include <map>
#include <string>
using namespace std;

void CREATE_TRL(vector<string>& trls, vector<string>& stops, map <string, vector<string>>& route);
void TRL_IN_STOP(map <string, vector<string>>& route);
void STOPS_IN_TRL(map <string, vector<string>>& route);
void TRLS(vector<string>& trls, map <string, vector<string>>& route);
