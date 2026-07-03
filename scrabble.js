// scrabble.js
#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');
const os = require('os');
const readline = require('readline');

const COLORS = {
    reset: '\x1b[0m',
    green: '\x1b[92m',
    yellow: '\x1b[93m',
    blue: '\x1b[94m',
    red: '\x1b[91m',
    cyan: '\x1b[96m',
    gray: '\x1b[90m',
    bold: '\x1b[1m'
};

function colorize(text, color) {
    return COLORS[color] + text + COLORS.reset;
}

const LETTER_SCORES = {
    'а':1,'б':3,'в':1,'г':3,'д':2,'е':1,'ё':3,'ж':5,'з':5,
    'и':1,'й':4,'к':2,'л':2,'м':2,'н':1,'о':1,'п':2,'р':1,
    'с':1,'т':1,'у':2,'ф':10,'х':5,'ц':5,'ч':5,'ш':8,'щ':10,
    'ъ':10,'ы':4,'ь':3,'э':8,'ю':8,'я':3
};

const LETTER_BAG = {
    'а':8,'б':2,'в':4,'г':2,'д':4,'е':9,'ё':1,'ж':1,'з':2,
    'и':6,'й':1,'к':4,'л':4,'м':3,'н':5,'о':10,'п':4,'р':5,
    'с':5,'т':5,'у':3,'ф':1,'х':1,'ц':1,'ч':1,'ш':1,'щ':1,
    'ъ':1,'ы':2,'ь':2,'э':1,'ю':1,'я':2
};

const DICT = new Set([
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
]);

class Scrabble {
    constructor(mode, dictPath) {
        this.mode = mode;
        this.dict = DICT;
        if (dictPath && fs.existsSync(dictPath)) {
            const content = fs.readFileSync(dictPath, 'utf8');
            content.split('\n').forEach(word => {
                const w = word.trim().toLowerCase();
                if (w.length >= 2) this.dict.add(w);
            });
        }
        this.bag = [];
        this.createBag();
        this.players = [
            { name: 'Игрок 1', letters: [], score: 0, passes: 0 },
            { name: this.mode === 'ai' ? 'Компьютер' : 'Игрок 2', letters: [], score: 0, passes: 0 }
        ];
        this.current = 0;
        this.gameOver = false;
        for (const p of this.players) {
            p.letters = this.drawLetters(7);
        }
    }

    createBag() {
        this.bag = [];
        for (const [ch, count] of Object.entries(LETTER_BAG)) {
            for (let i = 0; i < count; i++) {
                this.bag.push(ch);
            }
        }
        // shuffle
        for (let i = this.bag.length - 1; i > 0; i--) {
            const j = Math.floor(Math.random() * (i + 1));
            [this.bag[i], this.bag[j]] = [this.bag[j], this.bag[i]];
        }
    }

    drawLetters(n) {
        const drawn = [];
        for (let i = 0; i < n && this.bag.length > 0; i++) {
            drawn.push(this.bag.pop());
        }
        return drawn;
    }

    scoreWord(word) {
        let score = 0;
        for (const ch of word) {
            score += LETTER_SCORES[ch] || 0;
        }
        if (word.length >= 5) {
            score += (word.length - 4) * 5;
        }
        return score;
    }

    canMakeWord(word, letters) {
        const count = {};
        for (const ch of letters) count[ch] = (count[ch] || 0) + 1;
        for (const ch of word) {
            if (!count[ch]) return false;
            count[ch]--;
        }
        return true;
    }

    isValidWord(word) {
        return word.length >= 2 && this.dict.has(word);
    }

    findPossibleWords(letters) {
        const result = [];
        for (const w of this.dict) {
            if (w.length <= letters.length && this.canMakeWord(w, letters)) {
                result.push(w);
            }
        }
        return result;
    }

    aiMove() {
        const possible = this.findPossibleWords(this.players[1].letters);
        if (possible.length === 0) return null;
        let best = possible[0];
        for (const w of possible) {
            const sc = this.scoreWord(w);
            const bestSc = this.scoreWord(best);
            if (sc > bestSc || (sc === bestSc && w.length > best.length)) {
                best = w;
            }
        }
        return best;
    }

    display() {
        console.log(colorize('='.repeat(40), 'bold'));
        console.log(colorize('Скраббл (упрощённый)', 'bold'));
        console.log(`Счёт: ${this.players[0].name} = ${this.players[0].score}, ${this.players[1].name} = ${this.players[1].score}`);
        console.log(`Букв в мешке: ${this.bag.length}`);
        for (const p of this.players) {
            const letters = p.letters.map(ch => colorize(ch.toUpperCase(), 'cyan')).join(' ');
            console.log(`${p.name}: ${letters}`);
        }
    }

    playTurn(input) {
        const p = this.players[this.current];
        if (input === 'pass') {
            p.passes++;
            this.endTurn();
            return true;
        }
        const word = input.toLowerCase();
        if (!this.isValidWord(word)) {
            console.log(colorize('Слово не найдено или слишком короткое.', 'red'));
            return false;
        }
        if (!this.canMakeWord(word, p.letters)) {
            console.log(colorize('У вас нет нужных букв.', 'red'));
            return false;
        }
        const score = this.scoreWord(word);
        p.score += score;
        // remove letters
        for (const ch of word) {
            const idx = p.letters.indexOf(ch);
            if (idx !== -1) p.letters.splice(idx, 1);
        }
        const drawn = this.drawLetters(7 - p.letters.length);
        p.letters.push(...drawn);
        p.passes = 0;
        this.endTurn();
        console.log(colorize(`Слово '${word}' засчитано! +${score} очков.`, 'green'));
        return true;
    }

    endTurn() {
        this.current = 1 - this.current;
        if (this.mode === 'ai' && this.current === 1) {
            if (this.findPossibleWords(this.players[1].letters).length === 0 && this.bag.length === 0) {
                this.players[1].passes++;
                this.current = 0;
            }
        }
        if (this.players[0].passes >= 2 && this.players[1].passes >= 2) this.gameOver = true;
        if (this.bag.length === 0 && this.players[0].letters.length === 0 && this.players[1].letters.length === 0) {
            this.gameOver = true;
        }
    }

    async play() {
        const rl = readline.createInterface({
            input: process.stdin,
            output: process.stdout
        });
        const question = (q) => new Promise(resolve => rl.question(q, resolve));

        console.log(colorize('Добро пожаловать в упрощённый Скраббл!', 'bold'));
        console.log('Правила: составляйте слова из своих букв (минимум 2).');
        console.log('Очки = сумма стоимости букв + бонус за длину (5+ -> +5, 6+ -> +10, 7+ -> +15).');
        console.log("Введите 'pass' чтобы пропустить ход, '?' для подсказки, 'q' для выхода.\n");
        this.display();

        while (!this.gameOver) {
            const p = this.players[this.current];
            console.log(`\nХод: ${p.name}`);
            if (this.mode === 'ai' && this.current === 1) {
                const word = this.aiMove();
                if (!word) {
                    console.log('Компьютер не может составить слово. Пропуск хода.');
                    this.playTurn('pass');
                } else {
                    this.playTurn(word);
                    console.log(`Компьютер составил слово '${word}' (+${this.scoreWord(word)} очков)`);
                }
                this.display();
                continue;
            }
            while (true) {
                const input = await question('Ваш ход: ');
                if (input === 'q') {
                    console.log('Выход.');
                    rl.close();
                    return;
                }
                if (input === '?') {
                    const possible = this.findPossibleWords(p.letters);
                    if (possible.length) {
                        console.log('Возможные слова (первые 10):');
                        for (let i = 0; i < Math.min(10, possible.length); i++) {
                            console.log(`  ${possible[i]} (${this.scoreWord(possible[i])} очков)`);
                        }
                    } else {
                        console.log('Нет возможных слов.');
                    }
                    continue;
                }
                if (this.playTurn(input)) {
                    this.display();
                    break;
                }
            }
        }
        console.log('\nИгра окончена!');
        if (this.players[0].score > this.players[1].score) {
            console.log(`Победил ${this.players[0].name} со счётом ${this.players[0].score}!`);
        } else if (this.players[1].score > this.players[0].score) {
            console.log(`Победил ${this.players[1].name} со счётом ${this.players[1].score}!`);
        } else {
            console.log('Ничья!');
        }
        const best = Math.max(this.players[0].score, this.players[1].score);
        const recFile = path.join(os.homedir(), '.scrabble_record.json');
        let record = 0;
        try {
            const data = JSON.parse(fs.readFileSync(recFile, 'utf8'));
            record = data.record || 0;
        } catch {}
        if (best > record) {
            fs.writeFileSync(recFile, JSON.stringify({ record: best }));
            console.log(colorize(`Новый рекорд: ${best} очков!`, 'green'));
        } else {
            console.log(`Лучший рекорд: ${record} очков.`);
        }
        rl.close();
    }
}

async function main() {
    const args = process.argv.slice(2);
    let mode = 'vs';
    let dictPath = '';
    let showRecord = false;
    for (let i = 0; i < args.length; i++) {
        const arg = args[i];
        if (arg === 'ai') mode = 'ai';
        else if (arg === 'vs') mode = 'vs';
        else if (arg === '-r' || arg === '--record') showRecord = true;
        else if (arg === '-d' && i+1 < args.length) {
            dictPath = args[++i];
        } else if (arg === '-h' || arg === '--help') {
            console.log('Использование: node scrabble.js [vs|ai] [-r] [-d словарь]');
            return;
        }
    }
    if (showRecord) {
        const recFile = path.join(os.homedir(), '.scrabble_record.json');
        try {
            const data = JSON.parse(fs.readFileSync(recFile, 'utf8'));
            console.log(`Рекорд: ${data.record || 0} очков`);
        } catch {
            console.log('Рекордов пока нет.');
        }
        return;
    }
    const game = new Scrabble(mode, dictPath);
    await game.play();
}

main().catch(console.error);
