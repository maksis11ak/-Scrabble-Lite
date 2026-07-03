#!/usr/bin/env ruby
# scrabble.rb
# encoding: UTF-8

require 'json'
require 'fileutils'

COLORS = {
  reset: "\e[0m",
  green: "\e[92m",
  yellow: "\e[93m",
  blue: "\e[94m",
  red: "\e[91m",
  cyan: "\e[96m",
  gray: "\e[90m",
  bold: "\e[1m"
}

def colorize(text, color)
  "#{COLORS[color]}#{text}#{COLORS[:reset]}"
end

LETTER_SCORES = {
  'а'=>1,'б'=>3,'в'=>1,'г'=>3,'д'=>2,'е'=>1,'ё'=>3,'ж'=>5,'з'=>5,
  'и'=>1,'й'=>4,'к'=>2,'л'=>2,'м'=>2,'н'=>1,'о'=>1,'п'=>2,'р'=>1,
  'с'=>1,'т'=>1,'у'=>2,'ф'=>10,'х'=>5,'ц'=>5,'ч'=>5,'ш'=>8,'щ'=>10,
  'ъ'=>10,'ы'=>4,'ь'=>3,'э'=>8,'ю'=>8,'я'=>3
}

LETTER_BAG = {
  'а'=>8,'б'=>2,'в'=>4,'г'=>2,'д'=>4,'е'=>9,'ё'=>1,'ж'=>1,'з'=>2,
  'и'=>6,'й'=>1,'к'=>4,'л'=>4,'м'=>3,'н'=>5,'о'=>10,'п'=>4,'р'=>5,
  'с'=>5,'т'=>5,'у'=>3,'ф'=>1,'х'=>1,'ц'=>1,'ч'=>1,'ш'=>1,'щ'=>1,
  'ъ'=>1,'ы'=>2,'ь'=>2,'э'=>1,'ю'=>1,'я'=>2
}

DICT = Set.new(%w[
  абзац абонент автобус агрегат аквариум алгоритм
  амплитуда ананас анекдот антенна аппарат арбуз
  аромат артист архив аспект астроном атмосфера
  атом аудитория аэропорт базар баланс барабан
  бассейн батарея безопасность библиотека билет
  биология блокнот богатство болезнь бонус борщ
  ботинки брак бригада бронза буква бульвар бумага
  бутерброд быт бюджет вагон вариант вдохновение
  вектор вершина весна взаимодействие взгляд взрыв
  внимание воздух возраст война волонтер воображение
  воспитание впечатление время выбор выпуск выражение
  высота выступление газета галактика гарантия гармония
  гениальность география герой гитара глобус голос
  гора город государство грамота граница гриб
  груз гуманизм дар движение дворец дебют декада
  декорация делегат демократия деревня деталь диалог
  диплом директор дисциплина доброта договор дождь
  документ долг долина дом достижение достоинство
  драма друг дубль душа дым европа единство
  ежедневник желание железо жизнь журнал забота
  завод загадка закон зал запас запись защита
  звезда звук здание здоровье зеркало знак знание
  золото зона игра идея издание изображение
  изобретение интерес информация искусство история
  кабинет календарь камень канал капитал карьера
  катастрофа качество квартира кино класс климат
  книга ковер код количество коллектив команда
  комитет комната конкурс конструкция контакт
  контракт концерт копейка корень корзина корпус
  космос костюм кофе кран красота кредит кризис
  кристалл критерий круг крыша кубок культура
  курорт лаборатория лагерь ладонь лампа ландшафт
  лауреат лед лекция лес лето лечение лидер
  линия листок литература личность лоб ловушка
  логика локоть луч льгота любовь магазин магия
  макет максимум мальчик манеж маршрут масса
  математика материал матрица машина медицина мел
  мемориал меньшинство мера механизм микрофон миллион
  минута мир миссия мнение модель модернизация
  молоко момент монитор монумент море мост
  мотивация мощность музей музыка мышление навык
  нагрузка надежда название наличие народ наука
  находка нация небо неделя необходимость нефть
  низ новаторство норма ночь объект объем обучение
  общество объектив одежда озеро океан окно
  олимпиада операция опыт организация орден орел
  оригинал оркестр оружие осень основа ответ
  открытие отрасль отчет оценка память панорама
  парад парк пароль партия паспорт патриот пауза
  певец перемена период песня пианино письмо
  питание план планета пластик платформа племя
  пленум плоскость победа повод погода поддержка
  подход позиция познание показатель поколение поле
  полет политика половина пользователь помощь понятие
  порт портрет последствие постановка поток поэзия
  пояс правило практика предмет президент премия
  прибор приз приказ природа причина провинция
  прогноз программа продукт проект промышленность
  пропаганда проспект процесс процент профессия
  психология птица публика путь пьеса работа
  равновесие радио развитие размер разум район
  ранг расход реализация революция регион режиссер
  результат реклама рекомендация религия ремонт ресурс
  реформа рисунок ритм род роль роман рост
  рынок сад санкция сборник свет свобода связь
  сезон секрет сектор сельское хозяйство семья сервис
  серия сигнал сила символ система ситуация сказка
  скорость слава слово служба случай смысл событие
  совет сознание создание сок солнце соревнование
  состав состояние сотрудник сохранение союз спасение
  спектакль список спорт способ справедливость средство
  стабильность стандарт статья стекло стена степень
  стиль стол столица стоимость страна стратегия
  стремление строительство студент стул субъект судьба
  сумма сутки сцена счастье тайна талант танец
  театр текст телефон температура тенденция теория
  терапия термин территория техника технология тип
  тишина товар творчество темперамент темп течение
  транспорт требование третий труд туризм убеждение
  уважение уверенность удар удача узел указ
  украшение улица улучшение ум управление уровень
  урок успех установка устойчивость ученик учет
  фабрика факультет фигура физика философия фильм
  финал фирма флаг фокус фонд форма формула
  фотография фрагмент фронт функция характер химия
  хлеб хозяин холод хороший художник цвет цель
  центр цирк цифра часть человек черта чистота
  чувство шаг шанс школа шум экран эксперт
  экспорт элемент энергия эпизод эпоха эскиз этап
  эфир юбилей юмор юность яблоко явление язык
  январь яркость яхта
])

class Scrabble
  attr_reader :bag, :players, :current, :game_over, :mode, :dict

  def initialize(mode, dict_path=nil)
    @mode = mode
    @dict = DICT.dup
    if dict_path && File.exist?(dict_path)
      File.foreach(dict_path) do |line|
        w = line.strip.downcase
        @dict.add(w) if w.size >= 2
      end
    end
    @bag = []
    create_bag
    @players = [
      { name: "Игрок 1", letters: [], score: 0, passes: 0 },
      { name: mode == "ai" ? "Компьютер" : "Игрок 2", letters: [], score: 0, passes: 0 }
    ]
    @current = 0
    @game_over = false
    @players.each { |p| p[:letters] = draw_letters(7) }
  end

  def create_bag
    LETTER_BAG.each do |ch, count|
      count.times { @bag << ch }
    end
    @bag.shuffle!
  end

  def draw_letters(n)
    drawn = []
    n.times do
      break if @bag.empty?
      drawn << @bag.pop
    end
    drawn
  end

  def score_word(word)
    score = word.chars.sum { |ch| LETTER_SCORES[ch] || 0 }
    if word.size >= 5
      score += (word.size - 4) * 5
    end
    score
  end

  def can_make_word?(word, letters)
    cnt = Hash.new(0)
    letters.each { |ch| cnt[ch] += 1 }
    word.each_char do |ch|
      return false if cnt[ch].zero?
      cnt[ch] -= 1
    end
    true
  end

  def valid_word?(word)
    word.size >= 2 && @dict.include?(word)
  end

  def find_possible_words(letters)
    result = []
    @dict.each do |w|
      result << w if w.size <= letters.size && can_make_word?(w, letters)
    end
    result
  end

  def ai_move
    possible = find_possible_words(@players[1][:letters])
    return nil if possible.empty?
    best = possible.max_by { |w| [score_word(w), w.size] }
    best
  end

  def display
    puts colorize("=" * 40, :bold)
    puts colorize("Скраббл (упрощённый)", :bold)
    puts "Счёт: #{@players[0][:name]} = #{@players[0][:score]}, #{@players[1][:name]} = #{@players[1][:score]}"
    puts "Букв в мешке: #{@bag.size}"
    @players.each do |p|
      letters = p[:letters].map { |ch| colorize(ch.upcase, :cyan) }.join(' ')
      puts "#{p[:name]}: #{letters}"
    end
  end

  def play_turn(input)
    p = @players[@current]
    if input == "pass"
      p[:passes] += 1
      end_turn
      return true
    end
    word = input.downcase
    unless valid_word?(word)
      puts colorize("Слово не найдено или слишком короткое.", :red)
      return false
    end
    unless can_make_word?(word, p[:letters])
      puts colorize("У вас нет нужных букв.", :red)
      return false
    end
    score = score_word(word)
    p[:score] += score
    word.each_char do |ch|
      idx = p[:letters].index(ch)
      p[:letters].delete_at(idx) if idx
    end
    p[:letters].concat(draw_letters(7 - p[:letters].size))
    p[:passes] = 0
    end_turn
    puts colorize("Слово '#{word}' засчитано! +#{score} очков.", :green)
    true
  end

  def end_turn
    @current = 1 - @current
    if @mode == "ai" && @current == 1
      if find_possible_words(@players[1][:letters]).empty? && @bag.empty?
        @players[1][:passes] += 1
        @current = 0
      end
    end
    if @players[0][:passes] >= 2 && @players[1][:passes] >= 2
      @game_over = true
    end
    if @bag.empty? && @players[0][:letters].empty? && @players[1][:letters].empty?
      @game_over = true
    end
  end

  def play
    puts colorize("Добро пожаловать в упрощённый Скраббл!", :bold)
    puts "Правила: составляйте слова из своих букв (минимум 2)."
    puts "Очки = сумма стоимости букв + бонус за длину (5+ -> +5, 6+ -> +10, 7+ -> +15)."
    puts "Введите 'pass' чтобы пропустить ход, '?' для подсказки, 'q' для выхода.\n"
    display

    until @game_over
      p = @players[@current]
      puts "\nХод: #{p[:name]}"
      if @mode == "ai" && @current == 1
        word = ai_move
        if word.nil?
          puts "Компьютер не может составить слово. Пропуск хода."
          play_turn("pass")
        else
          play_turn(word)
          puts "Компьютер составил слово '#{word}' (+#{score_word(word)} очков)"
        end
        display
        next
      end
      loop do
        print "Ваш ход: "
        input = gets.chomp.strip
        if input == "q"
          puts "Выход."
          return
        end
        if input == "?"
          possible = find_possible_words(p[:letters])
          if possible.any?
            puts "Возможные слова (первые 10):"
            possible.first(10).each { |w| puts "  #{w} (#{score_word(w)} очков)" }
          else
            puts "Нет возможных слов."
          end
          next
        end
        if play_turn(input)
          display
          break
        end
      end
    end
    puts "\nИгра окончена!"
    if @players[0][:score] > @players[1][:score]
      puts "Победил #{@players[0][:name]} со счётом #{@players[0][:score]}!"
    elsif @players[1][:score] > @players[0][:score]
      puts "Победил #{@players[1][:name]} со счётом #{@players[1][:score]}!"
    else
      puts "Ничья!"
    end
    best = [@players[0][:score], @players[1][:score]].max
    rec_file = File.join(Dir.home, ".scrabble_record.json")
    record = 0
    if File.exist?(rec_file)
      begin
        data = JSON.parse(File.read(rec_file))
        record = data["record"] || 0
      rescue
      end
    end
    if best > record
      File.write(rec_file, JSON.pretty_generate({ "record" => best }))
      puts colorize("Новый рекорд: #{best} очков!", :green)
    else
      puts "Лучший рекорд: #{record} очков."
    end
  end
end

def main
  mode = "vs"
  dict_path = nil
  show_record = false
  ARGV.each_with_index do |arg, i|
    case arg
    when "ai" then mode = "ai"
    when "vs" then mode = "vs"
    when "-r", "--record" then show_record = true
    when "-d"
      dict_path = ARGV[i+1] if i+1 < ARGV.size
    when "-h", "--help"
      puts "Использование: ruby scrabble.rb [vs|ai] [-r] [-d словарь]"
      exit
    end
  end
  if show_record
    rec_file = File.join(Dir.home, ".scrabble_record.json")
    if File.exist?(rec_file)
      begin
        data = JSON.parse(File.read(rec_file))
        puts "Рекорд: #{data["record"] || 0} очков"
      rescue
        puts "Рекордов пока нет."
      end
    else
      puts "Рекордов пока нет."
    end
    return
  end
  game = Scrabble.new(mode, dict_path)
  game.play
end

main if __FILE__ == $0
