// scrabble.cpp
#include <iostream>
#include <vector>
#include <string>
#include <unordered_set>
#include <map>
#include <random>
#include <algorithm>
#include <fstream>
#include <cctype>
#include <filesystem>

using namespace std;
namespace fs = std::filesystem;

const string RESET = "\033[0m";
const string GREEN = "\033[92m";
const string YELLOW = "\033[93m";
const string BLUE = "\033[94m";
const string RED = "\033[91m";
const string CYAN = "\033[96m";
const string GRAY = "\033[90m";
const string BOLD = "\033[1m";

string colorize(const string& text, const string& color) {
    return color + text + RESET;
}

// Стоимость букв
unordered_map<char, int> LETTER_SCORES = {
    {'а',1},{'б',3},{'в',1},{'г',3},{'д',2},{'е',1},{'ё',3},{'ж',5},{'з',5},
    {'и',1},{'й',4},{'к',2},{'л',2},{'м',2},{'н',1},{'о',1},{'п',2},{'р',1},
    {'с',1},{'т',1},{'у',2},{'ф',10},{'х',5},{'ц',5},{'ч',5},{'ш',8},{'щ',10},
    {'ъ',10},{'ы',4},{'ь',3},{'э',8},{'ю',8},{'я',3}
};

// Мешок букв
unordered_map<char, int> LETTER_BAG = {
    {'а',8},{'б',2},{'в',4},{'г',2},{'д',4},{'е',9},{'ё',1},{'ж',1},{'з',2},
    {'и',6},{'й',1},{'к',4},{'л',4},{'м',3},{'н',5},{'о',10},{'п',4},{'р',5},
    {'с',5},{'т',5},{'у',3},{'ф',1},{'х',1},{'ц',1},{'ч',1},{'ш',1},{'щ',1},
    {'ъ',1},{'ы',2},{'ь',2},{'э',1},{'ю',1},{'я',2}
};

// Встроенный словарь (урезанный для примера)
unordered_set<string> DICT = {
    "абзац","абонент","автобус","агрегат","аквариум","алгоритм",
    "амплитуда","ананас","анекдот","антенна","аппарат","арбуз",
    "аромат","артист","архив","аспект","астроном","атмосфера",
    "атом","аудитория","аэропорт","базар","баланс","барабан",
    "бассейн","батарея","безопасность","библиотека","билет",
    "биология","блокнот","богатство","болезнь","бонус","борщ",
    "ботинки","брак","бригада","бронза","буква","бульвар","бумага",
    "бутерброд","быт","бюджет","вагон","вариант","вдохновение",
    "вектор","вершина","весна","взаимодействие","взгляд","взрыв",
    "внимание","воздух","возраст","война","волонтер","воображение",
    "воспитание","впечатление","время","выбор","выпуск","выражение",
    "высота","выступление","газета","галактика","гарантия","гармония",
    "гениальность","география","герой","гитара","глобус","голос",
    "гора","город","государство","грамота","граница","гриб","груз",
    "гуманизм","дар","движение","дворец","дебют","декада","декорация",
    "делегат","демократия","деревня","деталь","диалог","диплом",
    "директор","дисциплина","доброта","договор","дождь","документ",
    "долг","долина","дом","достижение","достоинство","драма","друг",
    "дубль","душа","дым","европа","единство","ежедневник","желание",
    "железо","жизнь","журнал","забота","завод","загадка","закон",
    "зал","запас","запись","защита","звезда","звук","здание","здоровье",
    "зеркало","знак","знание","золото","зона","игра","идея","издание",
    "изображение","изобретение","интерес","информация","искусство",
    "история","кабинет","календарь","камень","канал","капитал",
    "карьера","катастрофа","качество","квартира","кино","класс",
    "климат","книга","ковер","код","количество","коллектив","команда",
    "комитет","комната","конкурс","конструкция","контакт","контракт",
    "концерт","копейка","корень","корзина","корпус","космос","костюм",
    "кофе","кран","красота","кредит","кризис","кристалл","критерий",
    "круг","крыша","кубок","культура","курорт","лаборатория","лагерь",
    "ладонь","лампа","ландшафт","лауреат","лед","лекция","лес","лето",
    "лечение","лидер","линия","листок","литература","личность","лоб",
    "ловушка","логика","локоть","луч","льгота","любовь","магазин",
    "магия","макет","максимум","мальчик","манеж","маршрут","масса",
    "математика","материал","матрица","машина","медицина","мел",
    "мемориал","меньшинство","мера","механизм","микрофон","миллион",
    "минута","мир","миссия","мнение","модель","модернизация","молоко",
    "момент","монитор","монумент","море","мост","мотивация","мощность",
    "музей","музыка","мышление","навык","нагрузка","надежда","название",
    "наличие","народ","наука","находка","нация","небо","неделя",
    "необходимость","нефть","низ","новаторство","норма","ночь","объект",
    "объем","обучение","общество","объектив","одежда","озеро","океан",
    "окно","олимпиада","операция","опыт","организация","орден","орел",
    "оригинал","оркестр","оружие","осень","основа","ответ","открытие",
    "отрасль","отчет","оценка","память","панорама","парад","парк",
    "пароль","партия","паспорт","патриот","пауза","певец","перемена",
    "период","песня","пианино","письмо","питание","план","планета",
    "пластик","платформа","племя","пленум","плоскость","победа",
    "повод","погода","поддержка","подход","позиция","познание",
    "показатель","поколение","поле","полет","политика","половина",
    "пользователь","помощь","понятие","порт","портрет","последствие",
    "постановка","поток","поэзия","пояс","правило","практика","предмет",
    "президент","премия","прибор","приз","приказ","природа","причина",
    "провинция","прогноз","программа","продукт","проект","промышленность",
    "пропаганда","проспект","процесс","процент","профессия","психология",
    "птица","публика","путь","пьеса","работа","равновесие","радио",
    "развитие","размер","разум","район","ранг","расход","реализация",
    "революция","регион","режиссер","результат","реклама","рекомендация",
    "религия","ремонт","ресурс","реформа","рисунок","ритм","род",
    "роль","роман","рост","рынок","сад","санкция","сборник","свет",
    "свобода","связь","сезон","секрет","сектор","сельское хозяйство",
    "семья","сервис","серия","сигнал","сила","символ","система",
    "ситуация","сказка","скорость","слава","слово","служба","случай",
    "смысл","событие","совет","сознание","создание","сок","солнце",
    "соревнование","состав","состояние","сотрудник","сохранение",
    "союз","спасение","спектакль","список","спорт","способ","справедливость",
    "средство","стабильность","стандарт","статья","стекло","стена",
    "степень","стиль","стол","столица","стоимость","страна","стратегия",
    "стремление","строительство","студент","стул","субъект","судьба",
    "сумма","сутки","сцена","счастье","тайна","талант","танец",
    "театр","текст","телефон","температура","тенденция","теория",
    "терапия","термин","территория","техника","технология","тип",
    "тишина","товар","творчество","темперамент","темп","течение",
    "транспорт","требование","третий","труд","туризм","убеждение",
    "уважение","уверенность","удар","удача","узел","указ","украшение",
    "улица","улучшение","ум","управление","уровень","урок","успех",
    "установка","устойчивость","ученик","учет","фабрика","факультет",
    "фигура","физика","философия","фильм","финал","фирма","флаг",
    "фокус","фонд","форма","формула","фотография","фрагмент","фронт",
    "функция","характер","химия","хлеб","хозяин","холод","хороший",
    "художник","цвет","цель","центр","цирк","цифра","часть","человек",
    "черта","чистота","чувство","шаг","шанс","школа","шум","экран",
    "эксперт","экспорт","элемент","энергия","эпизод","эпоха","эскиз",
    "этап","эфир","юбилей","юмор","юность","яблоко","явление","язык",
    "январь","яркость","яхта"
};

struct Player {
    string name;
    vector<char> letters;
    int score = 0;
    int passes = 0;
};

class Scrabble {
public:
    Scrabble(string mode, string dictPath = "") {
        this->mode = mode;
        loadDict(dictPath);
        createBag();
        players[0].name = "Игрок 1";
        players[1].name = (mode == "ai") ? "Компьютер" : "Игрок 2";
        for (auto& p : players) {
            p.letters = drawLetters(7);
        }
        current = 0;
        gameOver = false;
    }

    void loadDict(const string& path) {
        if (!path.empty() && fs::exists(path)) {
            ifstream f(path);
            string word;
            while (f >> word) {
                if (word.size() >= 2) dict.insert(word);
            }
        } else {
            dict = DICT;
        }
    }

    void createBag() {
        bag.clear();
        for (auto& kv : LETTER_BAG) {
            for (int i=0; i<kv.second; ++i) bag.push_back(kv.first);
        }
        random_device rd;
        mt19937 g(rd());
        shuffle(bag.begin(), bag.end(), g);
    }

    vector<char> drawLetters(int n) {
        vector<char> drawn;
        for (int i=0; i<n && !bag.empty(); ++i) {
            drawn.push_back(bag.back());
            bag.pop_back();
        }
        return drawn;
    }

    int scoreWord(const string& word) {
        int s = 0;
        for (char ch : word) {
            auto it = LETTER_SCORES.find(ch);
            if (it != LETTER_SCORES.end()) s += it->second;
        }
        if (word.size() >= 5) s += (word.size() - 4) * 5;
        return s;
    }

    bool canMakeWord(const string& word, const vector<char>& letters) {
        multiset<char> ms(letters.begin(), letters.end());
        for (char ch : word) {
            auto it = ms.find(ch);
            if (it == ms.end()) return false;
            ms.erase(it);
        }
        return true;
    }

    bool isValidWord(const string& word) {
        return word.size() >= 2 && dict.find(word) != dict.end();
    }

    vector<string> findPossibleWords(const vector<char>& letters) {
        vector<string> result;
        for (const string& w : dict) {
            if (w.size() <= letters.size() && canMakeWord(w, letters)) {
                result.push_back(w);
            }
        }
        return result;
    }

    string aiMove() {
        auto possible = findPossibleWords(players[1].letters);
        if (possible.empty()) return "";
        string best = possible[0];
        for (auto& w : possible) {
            if (scoreWord(w) > scoreWord(best) || (scoreWord(w) == scoreWord(best) && w.size() > best.size())) {
                best = w;
            }
        }
        return best;
    }

    void display() {
        cout << colorize(string(40, '='), BOLD) << endl;
        cout << colorize("Скраббл (упрощённый)", BOLD) << endl;
        cout << "Счёт: " << players[0].name << " = " << players[0].score << ", "
             << players[1].name << " = " << players[1].score << endl;
        cout << "Букв в мешке: " << bag.size() << endl;
        for (int i=0; i<2; ++i) {
            string letters;
            for (char ch : players[i].letters) letters += colorize(string(1, toupper(ch)), CYAN) + " ";
            cout << players[i].name << ": " << letters << endl;
        }
    }

    bool playTurn(const string& input) {
        Player& p = players[current];
        if (input == "pass") {
            p.passes++;
            endTurn();
            return true;
        }
        string word = input;
        transform(word.begin(), word.end(), word.begin(), ::tolower);
        if (!isValidWord(word)) {
            cout << colorize("Слово не найдено или слишком короткое.", RED) << endl;
            return false;
        }
        if (!canMakeWord(word, p.letters)) {
            cout << colorize("У вас нет нужных букв.", RED) << endl;
            return false;
        }
        int score = scoreWord(word);
        p.score += score;
        // Удаляем использованные буквы
        for (char ch : word) {
            auto it = find(p.letters.begin(), p.letters.end(), ch);
            if (it != p.letters.end()) p.letters.erase(it);
        }
        // Добираем
        auto drawn = drawLetters(7 - p.letters.size());
        p.letters.insert(p.letters.end(), drawn.begin(), drawn.end());
        p.passes = 0;
        endTurn();
        cout << colorize("Слово '" + word + "' засчитано! +" + to_string(score) + " очков.", GREEN) << endl;
        return true;
    }

    void endTurn() {
        current = 1 - current;
        // Проверка окончания
        if (mode == "ai" && current == 1) {
            if (findPossibleWords(players[1].letters).empty() && bag.empty()) {
                players[1].passes++;
                current = 0;
            }
        }
        if (players[0].passes >= 2 && players[1].passes >= 2) gameOver = true;
        if (bag.empty() && players[0].letters.empty() && players[1].letters.empty()) gameOver = true;
    }

    bool isGameOver() { return gameOver; }

    string getWinner() {
        if (players[0].score > players[1].score) return players[0].name;
        else if (players[1].score > players[0].score) return players[1].name;
        else return "";
    }

    void play() {
        cout << colorize("Добро пожаловать в упрощённый Скраббл!", BOLD) << endl;
        cout << "Правила: составляйте слова из своих букв (минимум 2)." << endl;
        cout << "Очки = сумма стоимости букв + бонус за длину (5+ -> +5, 6+ -> +10, 7+ -> +15)." << endl;
        cout << "Введите 'pass' чтобы пропустить ход, '?' для подсказки, 'q' для выхода.\n" << endl;
        display();

        string input;
        while (!isGameOver()) {
            Player& p = players[current];
            cout << "\nХод: " << p.name << endl;
            if (mode == "ai" && current == 1) {
                string word = aiMove();
                if (word.empty()) {
                    cout << "Компьютер не может составить слово. Пропуск хода." << endl;
                    playTurn("pass");
                } else {
                    playTurn(word);
                    cout << "Компьютер составил слово '" << word << "' (+" << scoreWord(word) << " очков)" << endl;
                }
                display();
                continue;
            }
            while (true) {
                cout << "Ваш ход: ";
                getline(cin, input);
                if (input == "q") {
                    cout << "Выход." << endl;
                    return;
                }
                if (input == "?") {
                    auto possible = findPossibleWords(p.letters);
                    if (!possible.empty()) {
                        cout << "Возможные слова (первые 10):" << endl;
                        for (int i=0; i<min(10, (int)possible.size()); ++i) {
                            cout << "  " << possible[i] << " (" << scoreWord(possible[i]) << " очков)" << endl;
                        }
                    } else {
                        cout << "Нет возможных слов." << endl;
                    }
                    continue;
                }
                if (playTurn(input)) {
                    display();
                    break;
                }
            }
        }
        cout << "\nИгра окончена!" << endl;
        string winner = getWinner();
        if (!winner.empty()) {
            cout << "Победил " << winner << " со счётом " << max(players[0].score, players[1].score) << "!" << endl;
        } else {
            cout << "Ничья!" << endl;
        }
        // Рекорд
        int best = max(players[0].score, players[1].score);
        string home = getenv("HOME") ? getenv("HOME") : "";
        string recFile = home + "/.scrabble_record.json";
        int record = 0;
        if (fs::exists(recFile)) {
            ifstream f(recFile);
            string content((istreambuf_iterator<char>(f)), istreambuf_iterator<char>());
            size_t pos = content.find("\"record\":");
            if (pos != string::npos) {
                int val = stoi(content.substr(pos+9));
                record = val;
            }
        }
        if (best > record) {
            ofstream f(recFile);
            f << "{\"record\":" << best << "}";
            cout << colorize("Новый рекорд: " + to_string(best) + " очков!", GREEN) << endl;
        } else {
            cout << "Лучший рекорд: " << record << " очков." << endl;
        }
    }

private:
    string mode;
    vector<char> bag;
    unordered_set<string> dict;
    Player players[2];
    int current;
    bool gameOver;
};

int main(int argc, char* argv[]) {
    string mode = "vs";
    string dictPath = "";
    bool showRecord = false;
    for (int i=1; i<argc; ++i) {
        string arg = argv[i];
        if (arg == "ai") mode = "ai";
        else if (arg == "vs") mode = "vs";
        else if (arg == "-r" || arg == "--record") showRecord = true;
        else if (arg == "-d" && i+1 < argc) dictPath = argv[++i];
        else if (arg == "-h" || arg == "--help") {
            cout << "Использование: scrabble [vs|ai] [-r] [-d словарь]" << endl;
            return 0;
        }
    }
    if (showRecord) {
        string home = getenv("HOME") ? getenv("HOME") : "";
        string recFile = home + "/.scrabble_record.json";
        if (fs::exists(recFile)) {
            ifstream f(recFile);
            string content((istreambuf_iterator<char>(f)), istreambuf_iterator<char>());
            size_t pos = content.find("\"record\":");
            if (pos != string::npos) {
                int record = stoi(content.substr(pos+9));
                cout << "Рекорд: " << record << " очков" << endl;
            } else cout << "Рекордов пока нет." << endl;
        } else cout << "Рекордов пока нет." << endl;
        return 0;
    }
    Scrabble game(mode, dictPath);
    game.play();
    return 0;
}
