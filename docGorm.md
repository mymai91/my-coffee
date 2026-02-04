## 1️⃣ QUERY METHODS
### Find All

#### Rails:

```ruby
# Get all
Order.all

# Get by IDs
Order.find([1, 2, 3])
Order.where(id: [1, 2, 3])
```

#### GORM:

```go
// Get all
var orders []Order
db.Find(&orders)

// Get by IDs
db.Find(&orders, []int{1, 2, 3})
db.Where("id IN ?", []int{1, 2, 3}).Find(&orders)
```

### Find First/Last
#### Rails:

```ruby
# First
Order.first
Order.find(1)

# Last
Order.last

# First with condition
Order.where(status: 1).first
```

#### GORM:

```go
// First
var order Order
db.First(&order)
db.First(&order, 1)

// Last
db.Last(&order)

// First with condition
db.Where("status = ?", 1).First(&order)
```

### Where - Filtering
#### Rails:

```ruby
# Simple where
Order.where(status: 1)

# Multiple conditions
Order.where(status: 1, menu_item_name: "Espresso")
Order.where("status = ? AND menu_item_name = ?", 1, "Espresso")

# IN clause
Order.where(status: [1, 2, 3])

# LIKE
Order.where("menu_item_name LIKE ?", "%Latte%")

# NOT
Order.where.not(status: 1)

# OR
Order.where(status: 1).or(Order.where(status: 2))
```

#### GORM:

```go
// Simple where
db.Where("status = ?", 1).Find(&orders)
db.Where(&Order{Status: 1}).Find(&orders)

// Multiple conditions
db.Where("status = ? AND menu_item_name = ?", 1, "Espresso").Find(&orders)
db.Where(map[string]interface{}{"status": 1, "menu_item_name": "Espresso"}).Find(&orders)

// IN clause
db.Where("status IN ?", []int{1, 2, 3}).Find(&orders)

// LIKE
db.Where("menu_item_name LIKE ?", "%Latte%").Find(&orders)

// NOT
db.Not("status = ?", 1).Find(&orders)

// OR
db.Where("status = ?", 1).Or("status = ?", 2).Find(&orders)
```

### Select - Pick Columns
#### Rails:

```ruby
# Select specific columns
Order.select(:id, :menu_item_name)
Order.select("id, menu_item_name")

# With aggregate
Order.select("COUNT(*) as count")
```

#### GORM:

```go
// Select specific columns
db.Select("id, menu_item_name").Find(&orders)

// With aggregate
db.Select("COUNT(*) as count").Find(&result)
```

### Order - Sorting
#### Rails:

```ruby
# ASC
Order.order(:created_at)
Order.order("created_at ASC")

# DESC
Order.order(created_at: :desc)
Order.order("created_at DESC")

# Multiple
Order.order(:status, created_at: :desc)
```

#### GORM:

```go
// ASC
db.Order("created_at").Find(&orders)
db.Order("created_at ASC").Find(&orders)

// DESC
db.Order("created_at DESC").Find(&orders)

// Multiple
db.Order("status ASC, created_at DESC").Find(&orders)
```

### Limit & Offset - Pagination
#### Rails:

```ruby
# First 10
Order.limit(10)

# Skip 20, take 10
Order.offset(20).limit(10)

# Pagination
page = 2
per_page = 10
Order.offset((page - 1) * per_page).limit(per_page)
```

#### GORM:

```go
// First 10
db.Limit(10).Find(&orders)

// Skip 20, take 10
db.Offset(20).Limit(10).Find(&orders)

// Pagination
page := 2
pageSize := 10
db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders)
```

### Count
#### Rails:

```ruby
# Count all
Order.count

# Count with condition
Order.where(status: 1).count

# Count by group
Order.group(:status).count
# => {1 => 10, 2 => 5}
```

#### GORM:

```go
// Count all
var count int64
db.Model(&Order{}).Count(&count)

// Count with condition
db.Model(&Order{}).Where("status = ?", 1).Count(&count)

// Count by group
type Result struct {
    Status int
    Count  int64
}
var results []Result
db.Model(&Order{}).Select("status, COUNT(*) as count").Group("status").Scan(&results)
```

### Group & Having
#### Rails:

```ruby
# Group by
Order.group(:status).count

# With having
Order.select("status, COUNT(*) as count")
     .group(:status)
     .having("count > ?", 2)
```

#### GORM:

```go
// Group by
db.Select("status, COUNT(*) as count").Group("status").Find(&results)

// With having
db.Select("status, COUNT(*) as count")
  .Group("status")
  .Having("count > ?", 2)
  .Find(&results)
```

### Pluck - Get Single Column
#### Rails:

```ruby
# Pluck one column
Order.pluck(:menu_item_name)
# => ["Espresso", "Latte", "Cappuccino"]

# Pluck multiple columns
Order.pluck(:id, :menu_item_name)
# => [[1, "Espresso"], [2, "Latte"]]
```

#### GORM:

```go
// Pluck one column
var names []string
db.Model(&Order{}).Pluck("menu_item_name", &names)
// names = ["Espresso", "Latte", "Cappuccino"]

// Pluck multiple - use Scan
type Result struct {
    ID   uint
    Name string
}
var results []Result
db.Model(&Order{}).Select("id, menu_item_name as name").Scan(&results)
```

### Distinct
#### Rails:

```ruby
Order.distinct.pluck(:status)
```

#### GORM:

```go
var statuses []int
db.Model(&Order{}).Distinct("status").Pluck("status", &statuses)
```

---

## 2️⃣ CREATE METHODS

### Create
#### Rails:

```ruby
# Create one
Order.create(menu_item_name: "Espresso", status: 1)

# Create many
Order.create([
  { menu_item_name: "Espresso", status: 1 },
  { menu_item_name: "Latte", status: 1 }
])

# New + Save
order = Order.new(menu_item_name: "Espresso")
order.save
```

#### GORM:

```go
// Create one
order := Order{MenuItemName: "Espresso", Status: 1}
db.Create(&order)

// Create many
orders := []Order{
    {MenuItemName: "Espresso", Status: 1},
    {MenuItemName: "Latte", Status: 1},
}
db.Create(&orders)

// Manual save
order := Order{MenuItemName: "Espresso"}
db.Save(&order)
```

### Find or Create
#### Rails:

```ruby
# Find or create
Order.find_or_create_by(menu_item_name: "Espresso")

# Find or initialize (not save)
Order.find_or_initialize_by(menu_item_name: "Espresso")

# With defaults
Order.find_or_create_by(menu_item_name: "Espresso") do |order|
  order.status = 1
end
```

#### GORM:

```go
// Find or create
var order Order
db.Where(Order{MenuItemName: "Espresso"}).FirstOrCreate(&order)

// Find or initialize
db.Where(Order{MenuItemName: "Espresso"}).FirstOrInit(&order)

// With defaults
db.Where(Order{MenuItemName: "Espresso"})
  .Attrs(Order{Status: 1})
  .FirstOrCreate(&order)
```

---

## 3️⃣ UPDATE METHODS

### Update
#### Rails:

```ruby
# Update one field
order = Order.find(1)
order.update(status: 2)

# Update without loading
Order.where(id: 1).update(status: 2)

# Update multiple fields
order.update(status: 2, menu_item_name: "Latte")

# Update all matching
Order.where(status: 1).update_all(status: 2)
```

#### GORM:

```go
// Update one field
var order Order
db.First(&order, 1)
db.Model(&order).Update("status", 2)

// Update without loading
db.Model(&Order{}).Where("id = ?", 1).Update("status", 2)

// Update multiple fields
db.Model(&order).Updates(Order{Status: 2, MenuItemName: "Latte"})
db.Model(&order).Updates(map[string]interface{}{"status": 2, "menu_item_name": "Latte"})

// Update all matching
db.Model(&Order{}).Where("status = ?", 1).Update("status", 2)
```

### Increment/Decrement
#### Rails:

```ruby
order.increment!(:quantity)
order.decrement!(:quantity)
```

#### GORM:

```go
db.Model(&order).Update("quantity", gorm.Expr("quantity + ?", 1))
db.Model(&order).Update("quantity", gorm.Expr("quantity - ?", 1))
```

---

## 4️⃣ DELETE METHODS

### Delete
#### Rails:

```ruby
# Destroy one (with callbacks)
order = Order.find(1)
order.destroy

# Delete without callbacks
order.delete

# Destroy many
Order.where(status: 5).destroy_all

# Delete all
Order.delete_all
```

#### GORM:

```go
// Delete one (soft delete if has DeletedAt)
var order Order
db.First(&order, 1)
db.Delete(&order)

// Hard delete
db.Unscoped().Delete(&order)

// Delete many
db.Where("status = ?", 5).Delete(&Order{})

// Delete all
db.Where("1 = 1").Delete(&Order{})
```

---

## 5️⃣ ASSOCIATIONS & JOINS

### Setup Models
#### Rails:

```ruby
class User < ApplicationRecord
  has_many :orders
end

class Order < ApplicationRecord
  belongs_to :user
end
```

#### GORM:

```go
type User struct {
    gorm.Model
    Name   string
    Orders []Order  // has_many
}

type Order struct {
    gorm.Model
    UserID       uint
    MenuItemName string
    User         User  // belongs_to
}
```

### Preload (Eager Loading)
#### Rails:

```ruby
# Preload (N+1 solution)
Order.includes(:user)

# Access
orders.each do |order|
  puts order.user.name  # No extra query
end

# Nested preload
Order.includes(user: :profile)

# Preload with condition
User.includes(:orders).where(orders: { status: 1 })
```

#### GORM:

```go
// Preload (N+1 solution)
var orders []Order
db.Preload("User").Find(&orders)

// Access
for _, order := range orders {
    fmt.Println(order.User.Name)  // No extra query
}

// Nested preload
db.Preload("User.Profile").Find(&orders)

// Preload with condition
db.Preload("Orders", "status = ?", 1).Find(&users)
```

### Joins
#### Rails:

```ruby
# Inner join
Order.joins(:user)

# Join with condition
Order.joins(:user).where(users: { name: "John" })

# Left join
Order.left_joins(:user)

# Multiple joins
Order.joins(:user, :items)
```

#### GORM:

```go
// Inner join
db.Joins("User").Find(&orders)

// Join with condition
db.Joins("User").Where("users.name = ?", "John").Find(&orders)

// Left join
db.Joins("LEFT JOIN users ON users.id = orders.user_id").Find(&orders)

// Multiple joins
db.Joins("User").Joins("Items").Find(&orders)
```

### Real Example: 3 Tables
#### Rails:

```ruby
# Models
class User < ApplicationRecord
  has_many :orders
end

class Order < ApplicationRecord
  belongs_to :user
  has_many :order_items
end

class OrderItem < ApplicationRecord
  belongs_to :order
end

# Query: Load everything
users = User.includes(orders: :order_items)

users.each do |user|
  user.orders.each do |order|
    order.order_items.each do |item|
      puts item.menu_item_name
    end
  end
end

# Query: Join with conditions
Order.joins(:user)
     .joins(:order_items)
     .where(users: { email: /.*@gmail\.com/ })
     .where(orders: { status: 1 })
```

#### GORM:

```go
// Models
type User struct {
    gorm.Model
    Name   string
    Orders []Order
}

type Order struct {
    gorm.Model
    UserID uint
    User   User
    Items  []OrderItem
}

type OrderItem struct {
    gorm.Model
    OrderID      uint
    MenuItemName string
}

// Query: Load everything
var users []User
db.Preload("Orders.Items").Find(&users)

for _, user := range users {
    for _, order := range user.Orders {
        for _, item := range order.Items {
            fmt.Println(item.MenuItemName)
        }
    }
}

// Query: Join with conditions
var orders []Order
db.Joins("User").
   Joins("LEFT JOIN order_items ON order_items.order_id = orders.id").
   Where("users.email LIKE ?", "%@gmail.com").
   Where("orders.status = ?", 1).
   Find(&orders)
```

---

## 6️⃣ SCOPES
#### Rails:

```ruby
# Define scope
class Order < ApplicationRecord
  scope :queued, -> { where(status: 1) }
  scope :recent, -> { order(created_at: :desc).limit(10) }
end

# Use scope
Order.queued.recent
```

#### GORM:

```go
// Define scope
func Queued(db *gorm.DB) *gorm.DB {
    return db.Where("status = ?", 1)
}

func Recent(db *gorm.DB) *gorm.DB {
    return db.Order("created_at DESC").Limit(10)
}

// Use scope
db.Scopes(Queued, Recent).Find(&orders)
```

---

## 7️⃣ TRANSACTIONS
#### Rails:

```ruby
ActiveRecord::Base.transaction do
  order = Order.create!(menu_item_name: "Espresso")
  OrderItem.create!(order_id: order.id, name: "Milk")
end
# Auto rollback on error
```

#### GORM:

```go
db.Transaction(func(tx *gorm.DB) error {
    order := Order{MenuItemName: "Espresso"}
    if err := tx.Create(&order).Error; err != nil {
        return err  // Rollback
    }
    
    item := OrderItem{OrderID: order.ID, Name: "Milk"}
    if err := tx.Create(&item).Error; err != nil {
        return err  // Rollback
    }
    
    return nil  // Commit
})
```

---

## 8️⃣ RAW SQL
#### Rails:

```ruby
# Raw query
Order.find_by_sql("SELECT * FROM orders WHERE status = ?", 1)

# Execute
ActiveRecord::Base.connection.execute("UPDATE orders SET status = 2")
```

#### GORM:

```go
// Raw query
var orders []Order
db.Raw("SELECT * FROM orders WHERE status = ?", 1).Scan(&orders)

// Execute
db.Exec("UPDATE orders SET status = ?", 2)
```

---

## 9️⃣ CALLBACKS (Hooks)
#### Rails:

```ruby
class Order < ApplicationRecord
  before_create :set_default_status
  after_create :send_notification
  
  def set_default_status
    self.status ||= 1
  end
  
  def send_notification
    # Send email
  end
end
```

#### GORM:

```go
func (o *Order) BeforeCreate(tx *gorm.DB) error {
    if o.Status == 0 {
        o.Status = 1
    }
    return nil
}

func (o *Order) AfterCreate(tx *gorm.DB) error {
    // Send notification
    return nil
}
```

---

## 🔟 VALIDATION
#### Rails:

```ruby
class Order < ApplicationRecord
  validates :menu_item_name, presence: true
  validates :status, inclusion: { in: [1, 2, 3, 4, 5] }
end
```

#### GORM:

```go
// GORM doesn't have built-in validations
// Use another library or validate manually

func (o *Order) BeforeCreate(tx *gorm.DB) error {
    if o.MenuItemName == "" {
        return errors.New("menu_item_name is required")
    }
    if o.Status < 1 || o.Status > 5 {
        return errors.New("invalid status")
    }
    return nil
}
```

---

## 📊 Summary Comparison

| Feature | Rails | GORM |
|---------|-------|------|
| Find all | `Order.all` | `db.Find(&orders)` |
| Where | `Order.where(status: 1)` | `db.Where("status = ?", 1).Find(&orders)` |
| First | `Order.first` | `db.First(&order)` |
| Create | `Order.create(...)` | `db.Create(&order)` |
| Update | `order.update(status: 2)` | `db.Model(&order).Update("status", 2)` |
| Delete | `order.destroy` | `db.Delete(&order)` |
| Joins | `Order.joins(:user)` | `db.Joins("User").Find(&orders)` |
| Preload | `Order.includes(:user)` | `db.Preload("User").Find(&orders)` |
| Scopes | `scope :queued, -> { where(...) }` | `func Queued(db) *gorm.DB { return db.Where(...) }` |
| Callbacks | `before_create :method` | `func (o *Order) BeforeCreate(tx)` |

> **GORM ≈ 90% similar to ActiveRecord!** 😊


