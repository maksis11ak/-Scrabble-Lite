// scrabble.cs
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text.Json;
using System.Threading;

class ScrabbleGame
{
    static string Colorize(string text, string color)
    {
        string col = color switch
        {
            "green" => "\x1b[92m",
            "yellow" => "\x1b[93m",
            "blue" => "\x1b[94m",
            "red" => "\x1b[91m",
            "cyan" => "\x1b[96m",
            "gray" => "\x1b[90m",
            "bold" => "\x1b[1m",
            _ => "\x1b[0m"
        };
        return col + text + "\x1b[0m";
    }

    static Dictionary<char, int> LETTER_SCORES = new Dictionary<char, int> {
        {'а',1},{'б',3},{'в',1},{'г',3},{'д',2},{'е',1},{'ё',3},{'ж',5},{'з',5},
        {'и',1},{'й',4},{'к',2},{'л',2},{'м',2},{'н',1},{'о',1},{'п',2},{'р',1},
        {'с',1},{'т',1},{'у',2},{'ф',10},{'х',5},{'ц',5},{'ч',5},{'ш',8},{'щ',10},
        {'ъ',10},{'ы',4},{'ь',3},{'э',8},{'ю',8},{'я',3}
    };

    static Dictionary<char, int> LETTER_BAG = new Dictionary<char, int> {
        {'а',8},{'б',2},{'в',4},{'г',2},{'д',4},{'е',9},{'ё',1},{'ж',1},{'з',2},
        {'и',6},{'й',1},{'к',4},{'л',4},{'м',3},{'н',5},{'о',10},{'п',4},{'р',5},
        {'с',5},{'т',5},{'у',3},{'ф',1},{'х',1},{'ц',1},{'ч',1},{'ш',1},{'щ',1},
        {'ъ',1},{'ы',2},{'ь',2},{'э',1},{'ю',1},{'я',2}
    };

    static HashSet<string> DICT = new HashSet<string> {
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

    class Player
    {
        public string Name { get; set; }
        public List<char> Letters { get; set; }
        public int Score { get; set; }
        public int Passes { get; set; }
    }

    private string mode;
    private List<char> bag;
    private HashSet<string> dict;
    private Player[] players = new Player[2];
    private int current;
    private bool gameOver;

    public Scrabble(string mode, string dictPath)
    {
        this.mode = mode;
        dict = new HashSet<string>(DICT);
        if (!string.IsNullOrEmpty(dictPath) && File.Exists(dictPath))
        {
            foreach (var line in File.ReadLines(dictPath))
            {
                var w = line.Trim().ToLower();
                if (w.Length >= 2) dict.Add(w);
            }
        }
        CreateBag();
        players[0] = new Player { Name = "Игрок 1", Letters = new List<char>() };
        players[1] = new Player { Name = mode == "ai" ? "Компьютер" : "Игрок 2", Letters = new List<char>() };
        current = 0;
        gameOver = false;
        foreach (var p in players)
        {
            p.Letters.AddRange(DrawLetters(7));
        }
    }

    private void CreateBag()
    {
        bag = new List<char>();
        foreach (var kv in LETTER_BAG)
            for (int i = 0; i < kv.Value; i++)
                bag.Add(kv.Key);
        Random rnd = new Random();
        for (int i = bag.Count - 1; i > 0; i--)
        {
            int j = rnd.Next(i + 1);
            var tmp = bag[i];
            bag[i] = bag[j];
            bag[j] = tmp;
        }
    }

    private List<char> DrawLetters(int n)
    {
        var drawn = new List<char>();
        for (int i = 0; i < n && bag.Count > 0; i++)
        {
            drawn.Add(bag[bag.Count - 1]);
            bag.RemoveAt(bag.Count - 1);
        }
        return drawn;
    }

    private int ScoreWord(string word)
    {
        int score = 0;
        foreach (char ch in word)
            if (LETTER_SCORES.ContainsKey(ch)) score += LETTER_SCORES[ch];
        if (word.Length >= 5) score += (word.Length - 4) * 5;
        return score;
    }

    private bool CanMakeWord(string word, List<char> letters)
    {
        var count = new Dictionary<char, int>();
        foreach (var ch in letters)
            count[ch] = count.ContainsKey(ch) ? count[ch] + 1 : 1;
        foreach (char ch in word)
        {
            if (!count.ContainsKey(ch) || count[ch] == 0) return false;
            count[ch]--;
        }
        return true;
    }

    private bool IsValidWord(string word) => word.Length >= 2 && dict.Contains(word);

    private List<string> FindPossibleWords(List<char> letters)
    {
        var result = new List<string>();
        foreach (var w in dict)
        {
            if (w.Length <= letters.Count && CanMakeWord(w, letters))
                result.Add(w);
        }
        return result;
    }

    private string AiMove()
    {
        var possible = FindPossibleWords(players[1].Letters);
        if (possible.Count == 0) return null;
        string best = possible[0];
        foreach (var w in possible)
        {
            if (ScoreWord(w) > ScoreWord(best) || (ScoreWord(w) == ScoreWord(best) && w.Length > best.Length))
                best = w;
        }
        return best;
    }

    public void Display()
    {
        Console.WriteLine(Colorize(new string('=', 40), "bold"));
        Console.WriteLine(Colorize("Скраббл (упрощённый)", "bold"));
        Console.WriteLine($"Счёт: {players[0].Name} = {players[0].Score}, {players[1].Name} = {players[1].Score}");
        Console.WriteLine($"Букв в мешке: {bag.Count}");
        foreach (var p in players)
        {
            var letters = string.Join(" ", p.Letters.Select(ch => Colorize(ch.ToString().ToUpper(), "cyan")));
            Console.WriteLine($"{p.Name}: {letters}");
        }
    }

    public bool PlayTurn(string input)
    {
        var p = players[current];
        if (input == "pass")
        {
            p.Passes++;
            EndTurn();
            return true;
        }
        string word = input.ToLower();
        if (!IsValidWord(word))
        {
            Console.WriteLine(Colorize("Слово не найдено или слишком короткое.", "red"));
            return false;
        }
        if (!CanMakeWord(word, p.Letters))
        {
            Console.WriteLine(Colorize("У вас нет нужных букв.", "red"));
            return false;
        }
        int score = ScoreWord(word);
        p.Score += score;
        foreach (char ch in word)
            p.Letters.Remove(ch);
        p.Letters.AddRange(DrawLetters(7 - p.Letters.Count));
        p.Passes = 0;
        EndTurn();
        Console.WriteLine(Colorize($"Слово '{word}' засчитано! +{score} очков.", "green"));
        return true;
    }

    private void EndTurn()
    {
        current = 1 - current;
        if (mode == "ai" && current == 1)
        {
            if (FindPossibleWords(players[1].Letters).Count == 0 && bag.Count == 0)
            {
                players[1].Passes++;
                current = 0;
            }
        }
        if (players[0].Passes >= 2 && players[1].Passes >= 2) gameOver = true;
        if (bag.Count == 0 && players[0].Letters.Count == 0 && players[1].Letters.Count == 0)
            gameOver = true;
    }

    public void Play()
    {
        Console.WriteLine(Colorize("Добро пожаловать в упрощённый Скраббл!", "bold"));
        Console.WriteLine("Правила: составляйте слова из своих букв (минимум 2).");
        Console.WriteLine("Очки = сумма стоимости букв + бонус за длину (5+ -> +5, 6+ -> +10, 7+ -> +15).");
        Console.WriteLine("Введите 'pass' чтобы пропустить ход, '?' для подсказки, 'q' для выхода.\n");
        Display();

        while (!gameOver)
        {
            var p = players[current];
            Console.WriteLine($"\nХод: {p.Name}");
            if (mode == "ai" && current == 1)
            {
                string word = AiMove();
                if (word == null)
                {
                    Console.WriteLine("Компьютер не может составить слово. Пропуск хода.");
                    PlayTurn("pass");
                }
                else
                {
                    PlayTurn(word);
                    Console.WriteLine($"Компьютер составил слово '{word}' (+{ScoreWord(word)} очков)");
                }
                Display();
                continue;
            }
            while (true)
            {
                Console.Write("Ваш ход: ");
                string input = Console.ReadLine().Trim();
                if (input == "q")
                {
                    Console.WriteLine("Выход.");
                    return;
                }
                if (input == "?")
                {
                    var possible = FindPossibleWords(p.Letters);
                    if (possible.Count > 0)
                    {
                        Console.WriteLine("Возможные слова (первые 10):");
                        for (int i = 0; i < Math.Min(10, possible.Count); i++)
                            Console.WriteLine($"  {possible[i]} ({ScoreWord(possible[i])} очков)");
                    }
                    else Console.WriteLine("Нет возможных слов.");
                    continue;
                }
                if (PlayTurn(input))
                {
                    Display();
                    break;
                }
            }
        }
        Console.WriteLine("\nИгра окончена!");
        if (players[0].Score > players[1].Score)
            Console.WriteLine($"Победил {players[0].Name} со счётом {players[0].Score}!");
        else if (players[1].Score > players[0].Score)
            Console.WriteLine($"Победил {players[1].Name} со счётом {players[1].Score}!");
        else
            Console.WriteLine("Ничья!");

        int best = Math.Max(players[0].Score, players[1].Score);
        string recFile = Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.UserProfile), ".scrabble_record.json");
        int record = 0;
        if (File.Exists(recFile))
        {
            try
            {
                string json = File.ReadAllText(recFile);
                var data = JsonSerializer.Deserialize<Dictionary<string, int>>(json);
                record = data?.GetValueOrDefault("record", 0) ?? 0;
            }
            catch { }
        }
        if (best > record)
        {
            File.WriteAllText(recFile, JsonSerializer.Serialize(new Dictionary<string, int> { { "record", best } }));
            Console.WriteLine(Colorize($"Новый рекорд: {best} очков!", "green"));
        }
        else
            Console.WriteLine($"Лучший рекорд: {record} очков.");
    }

    static void Main(string[] args)
    {
        string mode = "vs", dictPath = null;
        bool showRecord = false;
        for (int i = 0; i < args.Length; i++)
        {
            if (args[i] == "ai") mode = "ai";
            else if (args[i] == "vs") mode = "vs";
            else if (args[i] == "-r" || args[i] == "--record") showRecord = true;
            else if (args[i] == "-d" && i + 1 < args.Length) dictPath = args[++i];
            else if (args[i] == "-h" || args[i] == "--help")
            {
                Console.WriteLine("Использование: scrabble [vs|ai] [-r] [-d словарь]");
                return;
            }
        }
        if (showRecord)
        {
            string recFile = Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.UserProfile), ".scrabble_record.json");
            if (File.Exists(recFile))
            {
                try
                {
                    string json = File.ReadAllText(recFile);
                    var data = JsonSerializer.Deserialize<Dictionary<string, int>>(json);
                    Console.WriteLine($"Рекорд: {data?.GetValueOrDefault("record", 0)} очков");
                }
                catch { Console.WriteLine("Рекордов пока нет."); }
            }
            else Console.WriteLine("Рекордов пока нет.");
            return;
        }
        ScrabbleGame game = new ScrabbleGame(mode, dictPath);
        game.Play();
    }
}
