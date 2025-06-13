from ai import call_gpt
import json

# File Path 
GROCERY_FILE ='groceries.json'

def main():
    while True: 
        action = menu()

        if action == "exit":
            print("Exiting program. Goodbye!")
            break

def menu():
    """
    Menu options for the pantry, this is used by the user to select what would like to do
    1. List of groceries 
    2. Add groceries 
    3. Recepies 
    """
    menu_items = {
        '1': 'List of groceries', 
        '2':'Add Groceries', 
        '3': 'Recepies',
        '4': 'Exit'
    }
    
    print("\nPantry Menu:")
    for key, value in menu_items.items():
        print(f"{key}: {value}")
        
    print("")
    user_selection = input("Select your option: ")

    # Menu controler 
    if user_selection == "1":
        groceries_list()
    elif user_selection == "2":
        add_groceries()
    elif user_selection == "3":
        show_recipes()
    elif user_selection == "4":
        return "exit"
    else: 
        print("Invalid option. Try again!")
    
# Load groceries from JSON file 
def load_groceries():
    with open(GROCERY_FILE, 'r') as f:
        return json.load(f)

# Save groceries to JSON file 
def save_groceries(groceries):
    with open(GROCERY_FILE, 'w') as f:
        json.dump(groceries, f, indent=4)

def groceries_list():
    while True: 
        groceries = load_groceries()

        # Return list of groceries if available 
        if not groceries:
            print("No groceries found.")
        else:
            for item in groceries:
                print(f"- {item['quantity']} {item['unit']} of {item['name']}")

        print("")

        # return to the menu 
        cmd = input("Type 'back' to return to the menu:")
        if cmd.lower() == 'back':
            break

def add_groceries():
    while True:
        groceries = load_groceries()

        # item entry 
        name = input("Enter item name: ")
        quantity = input("Enter quantity: ")
        unit = input("Enter unit (e.g., pieces, liters): ")

        groceries.append({"name": name, "quantity": quantity, "unit": unit})
        save_groceries(groceries)

        print(f"\nAdded {quantity} {unit} of {name}")

        print("")

        # return to the menu 
        cmd = input("Type 'back' to return to the menu:")
        if cmd.lower() == 'back':
            break

def show_recipes():
    while True: 
        print("Displaying recipes...")

        print("")
        
        # return to the menu 
        cmd = input("Type 'back' to return to the menu:")
        if cmd.lower() == 'back':
            break


if __name__ == "__main__":
    main()