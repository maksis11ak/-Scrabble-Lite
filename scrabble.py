# scrabble.py
#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys
import random
import json
import os
from pathlib import Path
from collections import Counter

# ANSI цвета
COLORS = {
    'reset': '\033[0m',
    'green': '\033[92m',
    'yellow': '\033[93m',
    'blue': '\033[94m',
    'red': '\033[91m',
    'cyan': '\033[96m',
    'gray': '\033[90m',
    'bold': '\033[1m'
}

def colorize(text, color):
    return f"{COLORS.get(color, '')}{text}{COLORS['reset']}"

# Стоимость букв (русский язык)
LETTER_SCORES = {
    'а': 1, 'б': 3, 'в': 1, 'г': 3, 'д': 2, 'е': 1, 'ё': 3, 'ж': 5, 'з': 5,
    'и': 1, 'й': 4, 'к': 2, 'л': 2, 'м': 2, 'н': 1, 'о': 1, 'п': 2, 'р': 1,
    'с': 1, 'т': 1, 'у': 2, 'ф': 10, 'х': 5, 'ц': 5, 'ч': 5, 'ш': 8, 'щ': 10,
    'ъ': 10, 'ы': 4, 'ь': 3, 'э': 8, 'ю': 8, 'я': 3
}

# Мешок букв (частотность для русского языка)
LETTER_BAG = {
    'а': 8, 'б': 2, 'в': 4, 'г': 2, 'д': 4, 'е': 9, 'ё': 1, 'ж': 1, 'з': 2,
    'и': 6, 'й': 1, 'к': 4, 'л': 4, 'м': 3, 'н': 5, 'о': 10, 'п': 4, 'р': 5,
    'с': 5, 'т': 5, 'у': 3, 'ф': 1, 'х': 1, 'ц': 1, 'ч': 1, 'ш': 1, 'щ': 1,
    'ъ': 1, 'ы': 2, 'ь': 2, 'э': 1, 'ю': 1, 'я': 2
}

# Встроенный словарь (сокращённый, для примера)
DICT_WORDS = [
    "абажур", "абзац", "абонент", "автобус", "агрегат", "аквариум", "аккорд",
    "аккумулятор", "алгоритм", "амплитуда", "ананас", "анекдот", "антенна",
    "аппарат", "арбуз", "аромат", "артист", "архив", "аспект", "астроном",
    "атмосфера", "атом", "аудитория", "аэропорт", "базар", "баланс", "барабан",
    "бассейн", "батарея", "безопасность", "библиотека", "билет", "биология",
    "блокнот", "богатство", "болезнь", "бонус", "борщ", "ботинки", "брак",
    "бригада", "бронза", "буква", "бульвар", "бумага", "бутерброд", "быт",
    "бюджет", "вагон", "вариант", "вдохновение", "вектор", "вершина", "весна",
    "взаимодействие", "взгляд", "взрыв", "внимание", "воздух", "возраст",
    "война", "волонтер", "воображение", "воспитание", "впечатление", "время",
    "выбор", "выпуск", "выражение", "высота", "выступление", "газета",
    "галактика", "гарантия", "гармония", "гениальность", "география", "герой",
    "гитара", "глобус", "голос", "гора", "город", "государство", "грамота",
    "граница", "гриб", "груз", "гуманизм", "дар", "движение", "дворец",
    "дебют", "декада", "декорация", "делегат", "демократия", "деревня",
    "деталь", "диалог", "диплом", "директор", "дисциплина", "доброта",
    "договор", "дождь", "документ", "долг", "долина", "дом", "достижение",
    "достоинство", "драма", "друг", "дубль", "душа", "дым", "европа",
    "единство", "ежедневник", "желание", "железо", "жизнь", "журнал",
    "забота", "завод", "загадка", "закон", "зал", "запас", "запись",
    "защита", "звезда", "звук", "здание", "здоровье", "зеркало", "знак",
    "знание", "золото", "зона", "игра", "идея", "издание", "изображение",
    "изобретение", "интерес", "информация", "искусство", "история", "кабинет",
    "календарь", "камень", "канал", "капитал", "карьера", "катастрофа",
    "качество", "квартира", "кино", "класс", "климат", "книга", "ковер",
    "код", "количество", "коллектив", "команда", "комитет", "комната",
    "конкурс", "конструкция", "контакт", "контракт", "концерт", "копейка",
    "корень", "корзина", "корпус", "космос", "костюм", "кофе", "кран",
    "красота", "кредит", "кризис", "кристалл", "критерий", "круг", "крыша",
    "кубок", "культура", "курорт", "лаборатория", "лагерь", "ладонь",
    "лампа", "ландшафт", "лауреат", "лед", "лекция", "лес", "лето",
    "лечение", "лидер", "линия", "листок", "литература", "личность", "лоб",
    "ловушка", "логика", "локоть", "луч", "льгота", "любовь", "магазин",
    "магия", "макет", "максимум", "мальчик", "манеж", "маршрут", "масса",
    "математика", "материал", "матрица", "машина", "медицина", "мел",
    "мемориал", "меньшинство", "мера", "механизм", "микрофон", "миллион",
    "минута", "мир", "миссия", "мнение", "модель", "модернизация", "молоко",
    "момент", "монитор", "монумент", "море", "мост", "мотивация", "мощность",
    "музей", "музыка", "мышление", "навык", "нагрузка", "надежда", "название",
    "наличие", "народ", "наука", "находка", "нация", "небо", "неделя",
    "необходимость", "нефть", "низ", "новаторство", "норма", "ночь", "объект",
    "объем", "обучение", "общество", "объектив", "одежда", "озеро", "океан",
    "окно", "олимпиада", "операция", "опыт", "организация", "орден", "орел",
    "оригинал", "оркестр", "оружие", "осень", "основа", "ответ", "открытие",
    "отрасль", "отчет", "оценка", "память", "панорама", "парад", "парк",
    "пароль", "партия", "паспорт", "патриот", "пауза", "певец", "перемена",
    "период", "песня", "пианино", "письмо", "питание", "план", "планета",
    "пластик", "платформа", "племя", "пленум", "плоскость", "победа",
    "повод", "погода", "поддержка", "подход", "позиция", "познание",
    "показатель", "поколение", "поле", "полет", "политика", "половина",
    "пользователь", "помощь", "понятие", "порт", "портрет", "последствие",
    "постановка", "поток", "поэзия", "пояс", "правило", "практика", "предмет",
    "президент", "премия", "прибор", "приз", "приказ", "природа", "причина",
    "провинция", "прогноз", "программа", "продукт", "проект", "промышленность",
    "пропаганда", "проспект", "процесс", "процент", "профессия", "психология",
    "птица", "публика", "путь", "пьеса", "работа", "равновесие", "радио",
    "развитие", "размер", "разум", "район", "ранг", "расход", "реализация",
    "революция", "регион", "режиссер", "результат", "реклама", "рекомендация",
    "религия", "ремонт", "ресурс", "реформа", "рисунок", "ритм", "род",
    "роль", "роман", "рост", "рынок", "сад", "санкция", "сборник", "свет",
    "свобода", "связь", "сезон", "секрет", "сектор", "сельское хозяйство",
    "семья", "сервис", "серия", "сигнал", "сила", "символ", "система",
    "ситуация", "сказка", "скорость", "слава", "слово", "служба", "случай",
    "смысл", "событие", "совет", "сознание", "создание", "сок", "солнце",
    "соревнование", "состав", "состояние", "сотрудник", "сохранение",
    "союз", "спасение", "спектакль", "список", "спорт", "способ", "справедливость",
    "средство", "стабильность", "стандарт", "статья", "стекло", "стена",
    "степень", "стиль", "стол", "столица", "стоимость", "страна", "стратегия",
    "стремление", "строительство", "студент", "стул", "субъект", "судьба",
    "сумма", "сутки", "сцена", "счастье", "тайна", "талант", "танец",
    "театр", "текст", "телефон", "температура", "тенденция", "теория",
    "терапия", "термин", "территория", "техника", "технология", "тип",
    "тишина", "товар", "творчество", "темперамент", "темп", "течение",
    "транспорт", "требование", "третий", "труд", "туризм", "убеждение",
    "уважение", "уверенность", "удар", "удача", "узел", "указ", "украшение",
    "улица", "улучшение", "ум", "управление", "уровень", "урок", "успех",
    "установка", "устойчивость", "ученик", "учет", "фабрика", "факультет",
    "фигура", "физика", "философия", "фильм", "финал", "фирма", "флаг",
    "фокус", "фонд", "форма", "формула", "фотография", "фрагмент", "фронт",
    "функция", "характер", "химия", "хлеб", "хозяин", "холод", "хороший",
    "художник", "цвет", "цель", "центр", "цирк", "цифра", "часть", "человек",
    "черта", "чистота", "чувство", "шаг", "шанс", "школа", "шум", "экран",
    "эксперт", "экспорт", "элемент", "энергия", "эпизод", "эпоха", "эскиз",
    "этап", "эфир", "юбилей", "юмор", "юность", "яблоко", "явление", "язык",
    "январь", "яркость", "яхта"
]
DICT_SET = set(DICT_WORDS)

class ScrabbleGame:
    def __init__(self, mode='vs', dict_path=None):
        self.mode = mode
        self.dict = DICT_SET
        if dict_path and os.path.exists(dict_path):
            with open(dict_path, 'r', encoding='utf-8') as f:
                self.dict = set(word.strip().lower() for word in f if len(word.strip()) >= 2)
        self.bag = self._create_bag()
        self.players = [{'name': 'Игрок 1', 'letters': [], 'score': 0, 'passes': 0},
                        {'name': 'Игрок 2' if mode == 'vs' else 'Компьютер', 'letters': [], 'score': 0, 'passes': 0}]
        self.current_player = 0
        self.turn = 0
        self.game_over = False
        # Раздача начальных букв
        for p in self.players:
            p['letters'] = self._draw_letters(7)

    def _create_bag(self):
        bag = []
        for letter, count in LETTER_BAG.items():
            bag.extend([letter] * count)
        random.shuffle(bag)
        return bag

    def _draw_letters(self, count):
        drawn = []
        for _ in range(min(count, len(self.bag))):
            drawn.append(self.bag.pop())
        return drawn

    def _score_word(self, word):
        score = sum(LETTER_SCORES.get(ch, 0) for ch in word)
        # Бонус за длину
        if len(word) >= 5:
            score += (len(word) - 4) * 5  # 5->+5, 6->+10, 7->+15
        return score

    def _can_make_word(self, word, letters):
        cnt = Counter(letters)
        for ch in word:
            if cnt.get(ch, 0) == 0:
                return False
            cnt[ch] -= 1
        return True

    def _is_valid_word(self, word):
        return len(word) >= 2 and word in self.dict

    def _find_possible_words(self, letters):
        # Находит все слова из словаря, которые можно составить из букв
        words = []
        # перебор всех слов словаря (неэффективно, но для демонстрации)
        for w in self.dict:
            if len(w) <= len(letters) and self._can_make_word(w, letters):
                words.append(w)
        return words

    def _ai_move(self):
        # ИИ ищет слово с максимальным счётом
        letters = self.players[1]['letters']
        possible = self._find_possible_words(letters)
        if not possible:
            return None
        best_word = max(possible, key=lambda w: (self._score_word(w), len(w)))
        return best_word

    def display(self):
        print(colorize("=" * 40, 'bold'))
        print(colorize("Скраббл (упрощённый)", 'bold'))
        print(f"Счёт: {self.players[0]['name']} = {self.players[0]['score']}, {self.players[1]['name']} = {self.players[1]['score']}")
        print(f"Букв в мешке: {len(self.bag)}")
        for i, p in enumerate(self.players):
            letters_str = ' '.join(colorize(ch.upper(), 'cyan') for ch in p['letters'])
            print(f"{p['name']}: {letters_str}")

    def play_turn(self, word):
        player = self.players[self.current_player]
        if word == 'pass':
            player['passes'] += 1
            self._end_turn()
            return True, "Ход пропущен."
        word = word.lower()
        if not self._is_valid_word(word):
            return False, "Слово не найдено в словаре или слишком короткое."
        if not self._can_make_word(word, player['letters']):
            return False, "У вас нет нужных букв."
        # Подсчёт очков
        score = self._score_word(word)
        player['score'] += score
        # Удаляем использованные буквы
        for ch in word:
            player['letters'].remove(ch)
        # Добираем буквы
        player['letters'].extend(self._draw_letters(7 - len(player['letters'])))
        player['passes'] = 0
        self._end_turn()
        return True, f"Слово '{word}' засчитано! +{score} очков."

    def _end_turn(self):
        self.current_player = 1 - self.current_player
        self.turn += 1
        # Проверка окончания игры
        if self.mode == 'ai' and self.current_player == 1:
            # Если компьютер не может ходить, он пропускает
            letters = self.players[1]['letters']
            if not self._find_possible_words(letters) and len(self.bag) == 0:
                self.players[1]['passes'] += 1
                self.current_player = 0
        # Если оба игрока пропустили ход подряд, игра окончена
        if self.players[0]['passes'] >= 2 and self.players[1]['passes'] >= 2:
            self.game_over = True
        # Если мешок пуст и у игрока нет букв, он не может ходить
        # Простая проверка: если у обоих нет букв и мешок пуст, игра окончена
        if len(self.bag) == 0 and all(len(p['letters']) == 0 for p in self.players):
            self.game_over = True

    def is_game_over(self):
        return self.game_over

    def get_winner(self):
        if self.players[0]['score'] > self.players[1]['score']:
            return self.players[0]['name']
        elif self.players[1]['score'] > self.players[0]['score']:
            return self.players[1]['name']
        else:
            return None

    def play(self):
        print(colorize("Добро пожаловать в упрощённый Скраббл!", 'bold'))
        print("Правила: составляйте слова из своих букв (минимум 2).")
        print("Очки = сумма стоимости букв + бонус за длину (5+ -> +5, 6+ -> +10, 7+ -> +15).")
        print("Введите 'pass' чтобы пропустить ход, '?' для подсказки, 'q' для выхода.\n")
        self.display()

        while not self.is_game_over():
            player = self.players[self.current_player]
            print(f"\nХод: {player['name']}")
            if self.mode == 'ai' and self.current_player == 1:
                # Ход компьютера
                word = self._ai_move()
                if word is None:
                    print("Компьютер не может составить слово. Пропуск хода.")
                    self.play_turn('pass')
                else:
                    self.play_turn(word)
                    print(f"Компьютер составил слово '{word.upper()}' (+{self._score_word(word)} очков)")
                self.display()
                continue
            # Ход человека
            while True:
                cmd = input("Ваш ход: ").strip()
                if cmd.lower() == 'q':
                    print("Выход.")
                    return
                if cmd == '?':
                    possible = self._find_possible_words(player['letters'])
                    if possible:
                        print("Возможные слова (первые 10):")
                        for w in possible[:10]:
                            print(f"  {w} ({self._score_word(w)} очков)")
                    else:
                        print("Нет возможных слов.")
                    continue
                ok, msg = self.play_turn(cmd)
                if ok:
                    print(msg)
                    self.display()
                    break
                else:
                    print(msg)

        # Конец игры
        print("\nИгра окончена!")
        winner = self.get_winner()
        if winner:
            print(f"Победил {winner} со счётом {max(p['score'] for p in self.players)}!")
        else:
            print("Ничья!")
        # Сохранение рекорда
        best = max(p['score'] for p in self.players)
        record_file = Path.home() / '.scrabble_record.json'
        record = 0
        if record_file.exists():
            with open(record_file, 'r') as f:
                data = json.load(f)
                record = data.get('record', 0)
        if best > record:
            with open(record_file, 'w') as f:
                json.dump({'record': best}, f)
            print(colorize(f"Новый рекорд: {best} очков!", 'green'))
        else:
            print(f"Лучший рекорд: {record} очков.")

def main():
    mode = 'vs'
    dict_path = None
    show_record = False
    if len(sys.argv) > 1:
        if sys.argv[1] == 'ai':
            mode = 'ai'
        elif sys.argv[1] == 'vs':
            mode = 'vs'
        elif sys.argv[1] == '-r' or sys.argv[1] == '--record':
            show_record = True
        else:
            print("Используйте: scrabble.py [vs|ai] [-r] [-d словарь]")
            sys.exit(1)
    for i, arg in enumerate(sys.argv):
        if arg == '-d' and i+1 < len(sys.argv):
            dict_path = sys.argv[i+1]
    if show_record:
        record_file = Path.home() / '.scrabble_record.json'
        if record_file.exists():
            with open(record_file, 'r') as f:
                data = json.load(f)
                print(f"Рекорд: {data.get('record', 0)} очков")
        else:
            print("Рекордов пока нет.")
        return
    game = ScrabbleGame(mode, dict_path)
    game.play()

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\nВыход.")
        sys.exit(0)
