import queue
import random
import threading
import time

from model import table_state, tables_list, menu

time_unit = 1

orders = queue.Queue()
orders.join()
orders_bundle = []

class Clients(threading.Thread):
    def __init__(self, *args, **kwargs):
        super(Clients, self).__init__(*args, **kwargs)

    def run(self):
        while True:
            time.sleep(1)
            self.create_order()

    @staticmethod
    def create_order():
        (table_id, table) = next(
            ((idx, table) for idx, table in enumerate(tables_list) if table['state'] == table_state[0]), (None, None))
        if table_id is not None:
            max_wait_time = 0
            food_choices = []
            for i in range(random.randint(1, 5)):
                choice = random.choice(menu)
                if max_wait_time < choice['preparation-time']:
                    max_wait_time = choice['preparation-time']
                food_choices.append(choice['id'])
            max_wait_time = max_wait_time * 1.3
            neworder_id = int(random.randint(1, 10000) * random.randint(1, 10000))
            neworder = {
                'table_id': table['id'],
                'id': neworder_id,
                'items': food_choices,
                'priority': random.randint(1, 5),
                'max_wait': max_wait_time,

            }
            orders.put(neworder)
            # orders_bundle to verify the order after receiving it from kitchen
            orders_bundle.append(neworder)

            tables_list[table_id]['state'] = table_state[1]
            tables_list[table_id]['order_id'] = neworder_id

        else:
            time.sleep(random.randint(2, 10) * time_unit)
            (table_id, table) = next(
                ((idx, table) for idx, table in enumerate(tables_list) if table['state'] == table_state[3]), (None, None))
            if table_id is not None:
                tables_list[table_id]['state'] = table_state[0]