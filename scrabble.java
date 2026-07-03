// scrabble.java
import java.io.*;
import java.nio.file.*;
import java.util.*;
import java.util.stream.*;

public class scrabble {
    private static final String RESET = "\u001B[0m";
    private static final String GREEN = "\u001B[92m";
    private static final String YELLOW = "\u001B[93m";
    private static final String BLUE = "\u001B[94m";
    private static final String RED = "\u001B[91m";
    private static final String CYAN = "\u001B[96m";
    private static final String GRAY = "\u001B[90m";
    private static final String BOLD = "\u001B[1m";

    private static String colorize(String text, String color) {
        return color + text + RESET;
    }

    private static final Map<Character, Integer> LETTER_SCORES = new HashMap<>();
    static {
        String letters = "абвгдеёжзийклмнопрстуфхцчшщъыьэюя";
        int[] scores = {1,3,1,3,2,1,3,5,5,1,4,2,2,2,1,1,2,1,1,1,2,10,5,5,5,8,10,10,4,3,8,8,3};
        for (int i=0; i<letters.length(); i++) {
            LETTER_SCORES.put(letters.charAt(i), scores[i]);
        }
    }

    private static final Map<Character, Integer> LETTER_BAG = new HashMap<>();
    static {
        String chars = "абвгдеёжзийклмнопрстуфхцчшщъыьэюя";
        int[] counts = {8,2,4,2,4,9,1,1,2,6,1,4,4,3,5,10,4,5,5,5,3,1,1,1,1,1,1,1,2,2,1,1,2};
        for (int i=0; i<chars.length(); i++) {
            LETTER_BAG.put(chars.charAt(i), counts[i]);
        }
    }

    private static Set<String> DICT = new HashSet<>(Arrays.asList(
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
    ));

    private static class Player {
        String name;
        List<Character> letters = new ArrayList<>();
        int score = 0;
        int passes = 0;
    }

    private String mode;
    private List<Character> bag = new ArrayList<>();
    private Set<String> dict;
    private Player[] players = new Player[2];
    private int current;
    private boolean gameOver;
    private BufferedReader reader = new BufferedReader(new InputStreamReader(System.in));

    public scrabble(String mode, String dictPath) throws IOException {
        this.mode = mode;
        dict = new HashSet<>(DICT);
        if (dictPath != null && Files.exists(Paths.get(dictPath))) {
            for (String line : Files.readAllLines(Paths.get(dictPath))) {
                String w = line.trim().toLowerCase();
                if (w.length() >= 2) dict.add(w);
            }
        }
        createBag();
        players[0] = new Player();
        players[0].name = "Игрок 1";
        players[1] = new Player();
        players[1].name = mode.equals("ai") ? "Компьютер" : "Игрок 2";
        for (Player p : players) {
            p.letters.addAll(drawLetters(7));
        }
        current = 0;
        gameOver = false;
    }

    private void createBag() {
        for (Map.Entry<Character, Integer> e : LETTER_BAG.entrySet()) {
            for (int i=0; i<e.getValue(); i++) {
                bag.add(e.getKey());
            }
        }
        Collections.shuffle(bag);
    }

    private List<Character> drawLetters(int n) {
        List<Character> drawn = new ArrayList<>();
        for (int i=0; i<n && !bag.isEmpty(); i++) {
            drawn.add(bag.remove(bag.size()-1));
        }
        return drawn;
    }

    private int scoreWord(String word) {
        int score = 0;
        for (char ch : word.toCharArray()) {
            score += LETTER_SCORES.getOrDefault(ch, 0);
        }
        if (word.length() >= 5) {
            score += (word.length() - 4) * 5;
        }
        return score;
    }

    private boolean canMakeWord(String word, List<Character> letters) {
        Map<Character, Integer> cnt = new HashMap<>();
        for (char ch : letters) cnt.put(ch, cnt.getOrDefault(ch, 0) + 1);
        for (char ch : word.toCharArray()) {
            if (!cnt.containsKey(ch) || cnt.get(ch) == 0) return false;
            cnt.put(ch, cnt.get(ch) - 1);
        }
        return true;
    }

    private boolean isValidWord(String word) {
        return word.length() >= 2 && dict.contains(word);
    }

    private List<String> findPossibleWords(List<Character> letters) {
        List<String> result = new ArrayList<>();
        for (String w : dict) {
            if (w.length() <= letters.size() && canMakeWord(w, letters)) {
                result.add(w);
            }
        }
        return result;
    }

    private String aiMove() {
        List<String> possible = findPossibleWords(players[1].letters);
        if (possible.isEmpty()) return null;
        String best = possible.get(0);
        for (String w : possible) {
            if (scoreWord(w) > scoreWord(best) || (scoreWord(w) == scoreWord(best) && w.length() > best.length())) {
                best = w;
            }
        }
        return best;
    }

    private void display() {
        System.out.println(colorize("=".repeat(40), BOLD));
        System.out.println(colorize("Скраббл (упрощённый)", BOLD));
        System.out.printf("Счёт: %s = %d, %s = %d\n", players[0].name, players[0].score, players[1].name, players[1].score);
        System.out.println("Букв в мешке: " + bag.size());
        for (Player p : players) {
            String letters = p.letters.stream()
                .map(ch -> colorize(Character.toString(ch).toUpperCase(), CYAN))
                .collect(Collectors.joining(" "));
            System.out.println(p.name + ": " + letters);
        }
    }

    private boolean playTurn(String input) {
        Player p = players[current];
        if (input.equals("pass")) {
            p.passes++;
            endTurn();
            return true;
        }
        String word = input.toLowerCase();
        if (!isValidWord(word)) {
            System.out.println(colorize("Слово не найдено или слишком короткое.", RED));
            return false;
        }
        if (!canMakeWord(word, p.letters)) {
            System.out.println(colorize("У вас нет нужных букв.", RED));
            return false;
        }
        int score = scoreWord(word);
        p.score += score;
        for (char ch : word.toCharArray()) {
            p.letters.remove(Character.valueOf(ch));
        }
        p.letters.addAll(drawLetters(7 - p.letters.size()));
        p.passes = 0;
        endTurn();
        System.out.println(colorize("Слово '" + word + "' засчитано! +" + score + " очков.", GREEN));
        return true;
    }

    private void endTurn() {
        current = 1 - current;
        if (mode.equals("ai") && current == 1) {
            if (findPossibleWords(players[1].letters).isEmpty() && bag.isEmpty()) {
                players[1].passes++;
                current = 0;
            }
        }
        if (players[0].passes >= 2 && players[1].passes >= 2) gameOver = true;
        if (bag.isEmpty() && players[0].letters.isEmpty() && players[1].letters.isEmpty()) {
            gameOver = true;
        }
    }

    public void play() throws IOException {
        System.out.println(colorize("Добро пожаловать в упрощённый Скраббл!", BOLD));
        System.out.println("Правила: составляйте слова из своих букв (минимум 2).");
        System.out.println("Очки = сумма стоимости букв + бонус за длину (5+ -> +5, 6+ -> +10, 7+ -> +15).");
        System.out.println("Введите 'pass' чтобы пропустить ход, '?' для подсказки, 'q' для выхода.\n");
        display();

        while (!gameOver) {
            Player p = players[current];
            System.out.println("\nХод: " + p.name);
            if (mode.equals("ai") && current == 1) {
                String word = aiMove();
                if (word == null) {
                    System.out.println("Компьютер не может составить слово. Пропуск хода.");
                    playTurn("pass");
                } else {
                    playTurn(word);
                    System.out.println("Компьютер составил слово '" + word + "' (+" + scoreWord(word) + " очков)");
                }
                display();
                continue;
            }
            while (true) {
                System.out.print("Ваш ход: ");
                String input = reader.readLine().trim();
                if (input.equals("q")) {
                    System.out.println("Выход.");
                    return;
                }
                if (input.equals("?")) {
                    List<String> possible = findPossibleWords(p.letters);
                    if (!possible.isEmpty()) {
                        System.out.println("Возможные слова (первые 10):");
                        for (int i=0; i<Math.min(10, possible.size()); i++) {
                            System.out.println("  " + possible.get(i) + " (" + scoreWord(possible.get(i)) + " очков)");
                        }
                    } else {
                        System.out.println("Нет возможных слов.");
                    }
                    continue;
                }
                if (playTurn(input)) {
                    display();
                    break;
                }
            }
        }
        System.out.println("\nИгра окончена!");
        if (players[0].score > players[1].score) {
            System.out.println("Победил " + players[0].name + " со счётом " + players[0].score + "!");
        } else if (players[1].score > players[0].score) {
            System.out.println("Победил " + players[1].name + " со счётом " + players[1].score + "!");
        } else {
            System.out.println("Ничья!");
        }

        int best = Math.max(players[0].score, players[1].score);
        String recFile = System.getProperty("user.home") + "/.scrabble_record.json";
        int record = 0;
        if (Files.exists(Paths.get(recFile))) {
            try {
                String content = new String(Files.readAllBytes(Paths.get(recFile)));
                // простой парсинг
                int idx = content.indexOf("\"record\":");
                if (idx != -1) {
                    int start = idx + 9;
                    int end = content.indexOf(',', start);
                    if (end == -1) end = content.indexOf('}', start);
                    record = Integer.parseInt(content.substring(start, end).trim());
                }
            } catch (Exception ignored) {}
        }
        if (best > record) {
            Files.write(Paths.get(recFile), ("{\"record\":" + best + "}").getBytes());
            System.out.println(colorize("Новый рекорд: " + best + " очков!", GREEN));
        } else {
            System.out.println("Лучший рекорд: " + record + " очков.");
        }
    }

    public static void main(String[] args) throws IOException {
        String mode = "vs";
        String dictPath = null;
        boolean showRecord = false;
        for (int i=0; i<args.length; i++) {
            String arg = args[i];
            if (arg.equals("ai")) mode = "ai";
            else if (arg.equals("vs")) mode = "vs";
            else if (arg.equals("-r") || arg.equals("--record")) showRecord = true;
            else if (arg.equals("-d") && i+1 < args.length) dictPath = args[++i];
            else if (arg.equals("-h") || arg.equals("--help")) {
                System.out.println("Использование: java scrabble [vs|ai] [-r] [-d словарь]");
                return;
            }
        }
        if (showRecord) {
            String recFile = System.getProperty("user.home") + "/.scrabble_record.json";
            if (Files.exists(Paths.get(recFile))) {
                try {
                    String content = new String(Files.readAllBytes(Paths.get(recFile)));
                    int idx = content.indexOf("\"record\":");
                    if (idx != -1) {
                        int start = idx + 9;
                        int end = content.indexOf(',', start);
                        if (end == -1) end = content.indexOf('}', start);
                        int record = Integer.parseInt(content.substring(start, end).trim());
                        System.out.println("Рекорд: " + record + " очков");
                    } else System.out.println("Рекордов пока нет.");
                } catch (Exception e) { System.out.println("Рекордов пока нет."); }
            } else System.out.println("Рекордов пока нет.");
            return;
        }
        scrabble game = new scrabble(mode, dictPath);
        game.play();
    }
}
