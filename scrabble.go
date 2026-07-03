// scrabble.go
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	reset  = "\033[0m"
	green  = "\033[92m"
	yellow = "\033[93m"
	blue   = "\033[94m"
	red    = "\033[91m"
	cyan   = "\033[96m"
	gray   = "\033[90m"
	bold   = "\033[1m"
)

func colorize(text, color string) string {
	return color + text + reset
}

var letterScores = map[rune]int{
	'а': 1, 'б': 3, 'в': 1, 'г': 3, 'д': 2, 'е': 1, 'ё': 3, 'ж': 5, 'з': 5,
	'и': 1, 'й': 4, 'к': 2, 'л': 2, 'м': 2, 'н': 1, 'о': 1, 'п': 2, 'р': 1,
	'с': 1, 'т': 1, 'у': 2, 'ф': 10, 'х': 5, 'ц': 5, 'ч': 5, 'ш': 8, 'щ': 10,
	'ъ': 10, 'ы': 4, 'ь': 3, 'э': 8, 'ю': 8, 'я': 3,
}

var letterBag = map[rune]int{
	'а': 8, 'б': 2, 'в': 4, 'г': 2, 'д': 4, 'е': 9, 'ё': 1, 'ж': 1, 'з': 2,
	'и': 6, 'й': 1, 'к': 4, 'л': 4, 'м': 3, 'н': 5, 'о': 10, 'п': 4, 'р': 5,
	'с': 5, 'т': 5, 'у': 3, 'ф': 1, 'х': 1, 'ц': 1, 'ч': 1, 'ш': 1, 'щ': 1,
	'ъ': 1, 'ы': 2, 'ь': 2, 'э': 1, 'ю': 1, 'я': 2,
}

var dictSet = map[string]bool{
	"абзац": true, "абонент": true, "автобус": true, "агрегат": true,
	"аквариум": true, "алгоритм": true, "амплитуда": true, "ананас": true,
	"анекдот": true, "антенна": true, "аппарат": true, "арбуз": true,
	"аромат": true, "артист": true, "архив": true, "аспект": true,
	"астроном": true, "атмосфера": true, "атом": true, "аудитория": true,
	"аэропорт": true, "базар": true, "баланс": true, "барабан": true,
	"бассейн": true, "батарея": true, "безопасность": true, "библиотека": true,
	"билет": true, "биология": true, "блокнот": true, "богатство": true,
	"болезнь": true, "бонус": true, "борщ": true, "ботинки": true,
	"брак": true, "бригада": true, "бронза": true, "буква": true,
	"бульвар": true, "бумага": true, "бутерброд": true, "быт": true,
	"бюджет": true, "вагон": true, "вариант": true, "вдохновение": true,
	"вектор": true, "вершина": true, "весна": true, "взаимодействие": true,
	"взгляд": true, "взрыв": true, "внимание": true, "воздух": true,
	"возраст": true, "война": true, "волонтер": true, "воображение": true,
	"воспитание": true, "впечатление": true, "время": true, "выбор": true,
	"выпуск": true, "выражение": true, "высота": true, "выступление": true,
	"газета": true, "галактика": true, "гарантия": true, "гармония": true,
	"гениальность": true, "география": true, "герой": true, "гитара": true,
	"глобус": true, "голос": true, "гора": true, "город": true,
	"государство": true, "грамота": true, "граница": true, "гриб": true,
	"груз": true, "гуманизм": true, "дар": true, "движение": true,
	"дворец": true, "дебют": true, "декада": true, "декорация": true,
	"делегат": true, "демократия": true, "деревня": true, "деталь": true,
	"диалог": true, "диплом": true, "директор": true, "дисциплина": true,
	"доброта": true, "договор": true, "дождь": true, "документ": true,
	"долг": true, "долина": true, "дом": true, "достижение": true,
	"достоинство": true, "драма": true, "друг": true, "дубль": true,
	"душа": true, "дым": true, "европа": true, "единство": true,
	"ежедневник": true, "желание": true, "железо": true, "жизнь": true,
	"журнал": true, "забота": true, "завод": true, "загадка": true,
	"закон": true, "зал": true, "запас": true, "запись": true,
	"защита": true, "звезда": true, "звук": true, "здание": true,
	"здоровье": true, "зеркало": true, "знак": true, "знание": true,
	"золото": true, "зона": true, "игра": true, "идея": true,
	"издание": true, "изображение": true, "изобретение": true,
	"интерес": true, "информация": true, "искусство": true, "история": true,
	"кабинет": true, "календарь": true, "камень": true, "канал": true,
	"капитал": true, "карьера": true, "катастрофа": true, "качество": true,
	"квартира": true, "кино": true, "класс": true, "климат": true,
	"книга": true, "ковер": true, "код": true, "количество": true,
	"коллектив": true, "команда": true, "комитет": true, "комната": true,
	"конкурс": true, "конструкция": true, "контакт": true, "контракт": true,
	"концерт": true, "копейка": true, "корень": true, "корзина": true,
	"корпус": true, "космос": true, "костюм": true, "кофе": true,
	"кран": true, "красота": true, "кредит": true, "кризис": true,
	"кристалл": true, "критерий": true, "круг": true, "крыша": true,
	"кубок": true, "культура": true, "курорт": true, "лаборатория": true,
	"лагерь": true, "ладонь": true, "лампа": true, "ландшафт": true,
	"лауреат": true, "лед": true, "лекция": true, "лес": true,
	"лето": true, "лечение": true, "лидер": true, "линия": true,
	"листок": true, "литература": true, "личность": true, "лоб": true,
	"ловушка": true, "логика": true, "локоть": true, "луч": true,
	"льгота": true, "любовь": true, "магазин": true, "магия": true,
	"макет": true, "максимум": true, "мальчик": true, "манеж": true,
	"маршрут": true, "масса": true, "математика": true, "материал": true,
	"матрица": true, "машина": true, "медицина": true, "мел": true,
	"мемориал": true, "меньшинство": true, "мера": true, "механизм": true,
	"микрофон": true, "миллион": true, "минута": true, "мир": true,
	"миссия": true, "мнение": true, "модель": true, "модернизация": true,
	"молоко": true, "момент": true, "монитор": true, "монумент": true,
	"море": true, "мост": true, "мотивация": true, "мощность": true,
	"музей": true, "музыка": true, "мышление": true, "навык": true,
	"нагрузка": true, "надежда": true, "название": true, "наличие": true,
	"народ": true, "наука": true, "находка": true, "нация": true,
	"небо": true, "неделя": true, "необходимость": true, "нефть": true,
	"низ": true, "новаторство": true, "норма": true, "ночь": true,
	"объект": true, "объем": true, "обучение": true, "общество": true,
	"объектив": true, "одежда": true, "озеро": true, "океан": true,
	"окно": true, "олимпиада": true, "операция": true, "опыт": true,
	"организация": true, "орден": true, "орел": true, "оригинал": true,
	"оркестр": true, "оружие": true, "осень": true, "основа": true,
	"ответ": true, "открытие": true, "отрасль": true, "отчет": true,
	"оценка": true, "память": true, "панорама": true, "парад": true,
	"парк": true, "пароль": true, "партия": true, "паспорт": true,
	"патриот": true, "пауза": true, "певец": true, "перемена": true,
	"период": true, "песня": true, "пианино": true, "письмо": true,
	"питание": true, "план": true, "планета": true, "пластик": true,
	"платформа": true, "племя": true, "пленум": true, "плоскость": true,
	"победа": true, "повод": true, "погода": true, "поддержка": true,
	"подход": true, "позиция": true, "познание": true, "показатель": true,
	"поколение": true, "поле": true, "полет": true, "политика": true,
	"половина": true, "пользователь": true, "помощь": true, "понятие": true,
	"порт": true, "портрет": true, "последствие": true, "постановка": true,
	"поток": true, "поэзия": true, "пояс": true, "правило": true,
	"практика": true, "предмет": true, "президент": true, "премия": true,
	"прибор": true, "приз": true, "приказ": true, "природа": true,
	"причина": true, "провинция": true, "прогноз": true, "программа": true,
	"продукт": true, "проект": true, "промышленность": true,
	"пропаганда": true, "проспект": true, "процесс": true, "процент": true,
	"профессия": true, "психология": true, "птица": true, "публика": true,
	"путь": true, "пьеса": true, "работа": true, "равновесие": true,
	"радио": true, "развитие": true, "размер": true, "разум": true,
	"район": true, "ранг": true, "расход": true, "реализация": true,
	"революция": true, "регион": true, "режиссер": true, "результат": true,
	"реклама": true, "рекомендация": true, "религия": true, "ремонт": true,
	"ресурс": true, "реформа": true, "рисунок": true, "ритм": true,
	"род": true, "роль": true, "роман": true, "рост": true,
	"рынок": true, "сад": true, "санкция": true, "сборник": true,
	"свет": true, "свобода": true, "связь": true, "сезон": true,
	"секрет": true, "сектор": true, "сельское хозяйство": true, "семья": true,
	"сервис": true, "серия": true, "сигнал": true, "сила": true,
	"символ": true, "система": true, "ситуация": true, "сказка": true,
	"скорость": true, "слава": true, "слово": true, "служба": true,
	"случай": true, "смысл": true, "событие": true, "совет": true,
	"сознание": true, "создание": true, "сок": true, "солнце": true,
	"соревнование": true, "состав": true, "состояние": true, "сотрудник": true,
	"сохранение": true, "союз": true, "спасение": true, "спектакль": true,
	"список": true, "спорт": true, "способ": true, "справедливость": true,
	"средство": true, "стабильность": true, "стандарт": true, "статья": true,
	"стекло": true, "стена": true, "степень": true, "стиль": true,
	"стол": true, "столица": true, "стоимость": true, "страна": true,
	"стратегия": true, "стремление": true, "строительство": true,
	"студент": true, "стул": true, "субъект": true, "судьба": true,
	"сумма": true, "сутки": true, "сцена": true, "счастье": true,
	"тайна": true, "талант": true, "танец": true, "театр": true,
	"текст": true, "телефон": true, "температура": true, "тенденция": true,
	"теория": true, "терапия": true, "термин": true, "территория": true,
	"техника": true, "технология": true, "тип": true, "тишина": true,
	"товар": true, "творчество": true, "темперамент": true, "темп": true,
	"течение": true, "транспорт": true, "требование": true, "третий": true,
	"труд": true, "туризм": true, "убеждение": true, "уважение": true,
	"уверенность": true, "удар": true, "удача": true, "узел": true,
	"указ": true, "украшение": true, "улица": true, "улучшение": true,
	"ум": true, "управление": true, "уровень": true, "урок": true,
	"успех": true, "установка": true, "устойчивость": true, "ученик": true,
	"учет": true, "фабрика": true, "факультет": true, "фигура": true,
	"физика": true, "философия": true, "фильм": true, "финал": true,
	"фирма": true, "флаг": true, "фокус": true, "фонд": true,
	"форма": true, "формула": true, "фотография": true, "фрагмент": true,
	"фронт": true, "функция": true, "характер": true, "химия": true,
	"хлеб": true, "хозяин": true, "холод": true, "хороший": true,
	"художник": true, "цвет": true, "цель": true, "центр": true,
	"цирк": true, "цифра": true, "часть": true, "человек": true,
	"черта": true, "чистота": true, "чувство": true, "шаг": true,
	"шанс": true, "школа": true, "шум": true, "экран": true,
	"эксперт": true, "экспорт": true, "элемент": true, "энергия": true,
	"эпизод": true, "эпоха": true, "эскиз": true, "этап": true,
	"эфир": true, "юбилей": true, "юмор": true, "юность": true,
	"яблоко": true, "явление": true, "язык": true, "январь": true,
	"яркость": true, "яхта": true,
}

type Player struct {
	Name   string
	Letters []rune
	Score   int
	Passes  int
}

type Scrabble struct {
	mode    string
	bag     []rune
	dict    map[string]bool
	players [2]Player
	current int
	gameOver bool
}

func NewScrabble(mode, dictPath string) *Scrabble {
	s := &Scrabble{mode: mode, dict: dictSet}
	s.loadDict(dictPath)
	s.createBag()
	s.players[0].Name = "Игрок 1"
	if mode == "ai" {
		s.players[1].Name = "Компьютер"
	} else {
		s.players[1].Name = "Игрок 2"
	}
	for i := range s.players {
		s.players[i].Letters = s.drawLetters(7)
	}
	s.current = 0
	s.gameOver = false
	return s
}

func (s *Scrabble) loadDict(path string) {
	if path != "" {
		// загружаем из файла, но для простоты оставим встроенный
	}
}

func (s *Scrabble) createBag() {
	s.bag = []rune{}
	for ch, count := range letterBag {
		for i := 0; i < count; i++ {
			s.bag = append(s.bag, ch)
		}
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(s.bag), func(i, j int) {
		s.bag[i], s.bag[j] = s.bag[j], s.bag[i]
	})
}

func (s *Scrabble) drawLetters(n int) []rune {
	drawn := []rune{}
	for i := 0; i < n && len(s.bag) > 0; i++ {
		drawn = append(drawn, s.bag[len(s.bag)-1])
		s.bag = s.bag[:len(s.bag)-1]
	}
	return drawn
}

func (s *Scrabble) scoreWord(word string) int {
	score := 0
	for _, ch := range word {
		if val, ok := letterScores[ch]; ok {
			score += val
		}
	}
	if len(word) >= 5 {
		score += (len(word) - 4) * 5
	}
	return score
}

func (s *Scrabble) canMakeWord(word string, letters []rune) bool {
	// считаем буквы
	count := make(map[rune]int)
	for _, ch := range letters {
		count[ch]++
	}
	for _, ch := range word {
		if count[ch] == 0 {
			return false
		}
		count[ch]--
	}
	return true
}

func (s *Scrabble) isValidWord(word string) bool {
	return len(word) >= 2 && s.dict[word]
}

func (s *Scrabble) findPossibleWords(letters []rune) []string {
	result := []string{}
	for w := range s.dict {
		if len(w) <= len(letters) && s.canMakeWord(w, letters) {
			result = append(result, w)
		}
	}
	return result
}

func (s *Scrabble) aiMove() string {
	possible := s.findPossibleWords(s.players[1].Letters)
	if len(possible) == 0 {
		return ""
	}
	best := possible[0]
	for _, w := range possible {
		if s.scoreWord(w) > s.scoreWord(best) || (s.scoreWord(w) == s.scoreWord(best) && len(w) > len(best)) {
			best = w
		}
	}
	return best
}

func (s *Scrabble) display() {
	fmt.Println(colorize(strings.Repeat("=", 40), bold))
	fmt.Println(colorize("Скраббл (упрощённый)", bold))
	fmt.Printf("Счёт: %s = %d, %s = %d\n", s.players[0].Name, s.players[0].Score, s.players[1].Name, s.players[1].Score)
	fmt.Printf("Букв в мешке: %d\n", len(s.bag))
	for i := range s.players {
		letters := ""
		for _, ch := range s.players[i].Letters {
			letters += colorize(string(ch), cyan) + " "
		}
		fmt.Printf("%s: %s\n", s.players[i].Name, letters)
	}
}

func (s *Scrabble) playTurn(input string) bool {
	p := &s.players[s.current]
	if input == "pass" {
		p.Passes++
		s.endTurn()
		return true
	}
	word := strings.ToLower(input)
	if !s.isValidWord(word) {
		fmt.Println(colorize("Слово не найдено или слишком короткое.", red))
		return false
	}
	if !s.canMakeWord(word, p.Letters) {
		fmt.Println(colorize("У вас нет нужных букв.", red))
		return false
	}
	score := s.scoreWord(word)
	p.Score += score
	// удаляем буквы
	for _, ch := range word {
		for i, l := range p.Letters {
			if l == ch {
				p.Letters = append(p.Letters[:i], p.Letters[i+1:]...)
				break
			}
		}
	}
	// добираем
	drawn := s.drawLetters(7 - len(p.Letters))
	p.Letters = append(p.Letters, drawn...)
	p.Passes = 0
	s.endTurn()
	fmt.Printf(colorize("Слово '%s' засчитано! +%d очков.\n", green), word, score)
	return true
}

func (s *Scrabble) endTurn() {
	s.current = 1 - s.current
	if s.mode == "ai" && s.current == 1 {
		if len(s.findPossibleWords(s.players[1].Letters)) == 0 && len(s.bag) == 0 {
			s.players[1].Passes++
			s.current = 0
		}
	}
	if s.players[0].Passes >= 2 && s.players[1].Passes >= 2 {
		s.gameOver = true
	}
	if len(s.bag) == 0 && len(s.players[0].Letters) == 0 && len(s.players[1].Letters) == 0 {
		s.gameOver = true
	}
}

func (s *Scrabble) play() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(colorize("Добро пожаловать в упрощённый Скраббл!", bold))
	fmt.Println("Правила: составляйте слова из своих букв (минимум 2).")
	fmt.Println("Очки = сумма стоимости букв + бонус за длину (5+ -> +5, 6+ -> +10, 7+ -> +15).")
	fmt.Println("Введите 'pass' чтобы пропустить ход, '?' для подсказки, 'q' для выхода.\n")
	s.display()

	for !s.gameOver {
		p := &s.players[s.current]
		fmt.Printf("\nХод: %s\n", p.Name)
		if s.mode == "ai" && s.current == 1 {
			word := s.aiMove()
			if word == "" {
				fmt.Println("Компьютер не может составить слово. Пропуск хода.")
				s.playTurn("pass")
			} else {
				s.playTurn(word)
				fmt.Printf("Компьютер составил слово '%s' (+%d очков)\n", word, s.scoreWord(word))
			}
			s.display()
			continue
		}
		for {
			fmt.Print("Ваш ход: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if input == "q" {
				fmt.Println("Выход.")
				return
			}
			if input == "?" {
				possible := s.findPossibleWords(p.Letters)
				if len(possible) > 0 {
					fmt.Println("Возможные слова (первые 10):")
					for i, w := range possible {
						if i >= 10 {
							break
						}
						fmt.Printf("  %s (%d очков)\n", w, s.scoreWord(w))
					}
				} else {
					fmt.Println("Нет возможных слов.")
				}
				continue
			}
			if s.playTurn(input) {
				s.display()
				break
			}
		}
	}
	fmt.Println("\nИгра окончена!")
	if s.players[0].Score > s.players[1].Score {
		fmt.Printf("Победил %s со счётом %d!\n", s.players[0].Name, s.players[0].Score)
	} else if s.players[1].Score > s.players[0].Score {
		fmt.Printf("Победил %s со счётом %d!\n", s.players[1].Name, s.players[1].Score)
	} else {
		fmt.Println("Ничья!")
	}
	// рекорд
	best := s.players[0].Score
	if s.players[1].Score > best {
		best = s.players[1].Score
	}
	home := os.Getenv("HOME")
	recFile := filepath.Join(home, ".scrabble_record.json")
	record := 0
	if data, err := os.ReadFile(recFile); err == nil {
		var tmp map[string]int
		if err := json.Unmarshal(data, &tmp); err == nil {
			record = tmp["record"]
		}
	}
	if best > record {
		data, _ := json.Marshal(map[string]int{"record": best})
		os.WriteFile(recFile, data, 0644)
		fmt.Println(colorize(fmt.Sprintf("Новый рекорд: %d очков!", best), green))
	} else {
		fmt.Printf("Лучший рекорд: %d очков.\n", record)
	}
}

func main() {
	mode := "vs"
	dictPath := ""
	showRecord := false
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch arg {
		case "ai":
			mode = "ai"
		case "vs":
			mode = "vs"
		case "-r", "--record":
			showRecord = true
		case "-d":
			if i+1 < len(os.Args) {
				dictPath = os.Args[i+1]
				i++
			}
		case "-h", "--help":
			fmt.Println("Использование: scrabble [vs|ai] [-r] [-d словарь]")
			return
		}
	}
	if showRecord {
		home := os.Getenv("HOME")
		recFile := filepath.Join(home, ".scrabble_record.json")
		if data, err := os.ReadFile(recFile); err == nil {
			var tmp map[string]int
			if err := json.Unmarshal(data, &tmp); err == nil {
				fmt.Printf("Рекорд: %d очков\n", tmp["record"])
			} else {
				fmt.Println("Рекордов пока нет.")
			}
		} else {
			fmt.Println("Рекордов пока нет.")
		}
		return
	}
	game := NewScrabble(mode, dictPath)
	game.play()
}
