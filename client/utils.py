import sys
import unittest

class Storage:
    _instance = None

    def __new__(cls):
        if cls._instance is None:
            cls._instance = super(Storage, cls).__new__(cls)
            cls._instance.data = {}
        return cls._instance

    def set(self, key, value):
        self.data[key] = value

    def get(self, key):
        return self.data.get(key)

    def delete(self, key):
        if key in self.data:
            del self.data[key]

    def get_all(self):
        return self.data


class QuietTestResult(unittest.TextTestResult):
    def addError(self, test, err):
        # Подавляем вывод traceback для ошибок
        self.errors.append((test, None))

    def addFailure(self, test, err):
        # Подавляем вывод traceback для провалов
        self.failures.append((test, None))

    def printErrors(self):
        # Не выводим ошибки
        pass


class QuietTestRunner(unittest.TextTestRunner):
    def _makeResult(self):
        result = QuietTestResult(self.stream, self.descriptions, self.verbosity)
        result.showAll = False  # Отключить подробный вывод
        result.dots = False     # Отключить точки (.)
        return result

    def run(self, test):
        # Переопределяем вывод
        with open('/dev/null', 'w') as devnull:
            sys.stderr = devnull  # Подавляем вывод ошибок
            # sys.stdout = devnull  # Подавляем стандартный вывод
            result = super().run(test)
            sys.stderr = sys.__stderr__  # Восстанавливаем стандартные потоки
            sys.stdout = sys.__stdout__
        return result