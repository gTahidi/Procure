package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

const DATABASE_NAME = "./procurement.db"

func SetupDatabaseSchema() {
	db, err := sql.Open("sqlite3", DATABASE_NAME+"?_foreign_keys=on") // Enable FK enforcement
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	log.Println("Successfully connected to SQLite database:", DATABASE_NAME)

	tables := []string{
		`CREATE TABLE IF NOT EXISTS Users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT UNIQUE NOT NULL, email TEXT UNIQUE NOT NULL,
            role TEXT NOT NULL CHECK(role IN ('admin', 'procurement_officer', 'requester', 'supplier', 'approver', 'evaluator')),
            department TEXT, contactNumber TEXT, isActive INTEGER DEFAULT 1 CHECK(isActive IN (0,1))
        );`,
		`CREATE TABLE IF NOT EXISTS Suppliers (
            id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT UNIQUE NOT NULL, contact_person TEXT,
            email TEXT UNIQUE, phone TEXT
        );`,
		`CREATE TABLE IF NOT EXISTS Assets (
            id INTEGER PRIMARY KEY AUTOINCREMENT, amrId TEXT UNIQUE NOT NULL, emrId TEXT UNIQUE, description TEXT,
            capitalizedValue REAL, depreciation_method TEXT, depreciation_usefulLife INTEGER,
            depreciation_annualRate REAL, proc_prId INTEGER, proc_poId INTEGER, proc_grnId INTEGER,
            FOREIGN KEY (proc_prId) REFERENCES Requisitions(id) ON DELETE SET NULL,
            FOREIGN KEY (proc_poId) REFERENCES PurchaseOrders(id) ON DELETE SET NULL,
            FOREIGN KEY (proc_grnId) REFERENCES Deliveries(id) ON DELETE SET NULL
        );`,
		`CREATE TABLE IF NOT EXISTS Requisitions (
            id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER NOT NULL, 
            type TEXT NOT NULL CHECK(type IN ('goods', 'services', 'fixed_asset')),
            aac TEXT CHECK(aac IN ('A', 'F', 'P')), materialGroup TEXT, exchangeRate REAL,
            status TEXT DEFAULT 'pending', created_at TEXT DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES Users(id)
        );`,
		`CREATE TABLE IF NOT EXISTS RequisitionItems (
            id INTEGER PRIMARY KEY AUTOINCREMENT, requisition_id INTEGER NOT NULL,
            description TEXT NOT NULL, quantity REAL NOT NULL, unit TEXT, estimated_unit_price REAL,
            freight_cost REAL DEFAULT 0, insurance_cost REAL DEFAULT 0, installation_cost REAL DEFAULT 0,
            amr_id INTEGER, value REAL,
            FOREIGN KEY (requisition_id) REFERENCES Requisitions(id) ON DELETE CASCADE,
            FOREIGN KEY (amr_id) REFERENCES Assets(id) ON DELETE SET NULL
        );`,
		`CREATE TABLE IF NOT EXISTS Tenders (
            id INTEGER PRIMARY KEY AUTOINCREMENT, requisition_id INTEGER UNIQUE,
            title TEXT NOT NULL, description TEXT, status TEXT DEFAULT 'draft',
            publish_date TEXT, closing_date TEXT,
            metadata_businessArea TEXT, metadata_costCategory TEXT, metadata_fundAvailabilityType TEXT,
            created_at TEXT DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (requisition_id) REFERENCES Requisitions(id) ON DELETE SET NULL
        );`,
		`CREATE TABLE IF NOT EXISTS TenderEvaluationCriteria (
            id INTEGER PRIMARY KEY AUTOINCREMENT, tender_id INTEGER NOT NULL,
            type TEXT NOT NULL CHECK(type IN ('technical', 'commercial', 'delivery', 'compliance')),
            criterion_text TEXT NOT NULL, weight REAL NOT NULL CHECK(weight >= 0 AND weight <= 100),
            sub_criteria_name TEXT, sub_criteria_max_score INTEGER, is_mandatory INTEGER DEFAULT 0 CHECK(is_mandatory IN (0,1)),
            FOREIGN KEY (tender_id) REFERENCES Tenders(id) ON DELETE CASCADE,
            UNIQUE (tender_id, type, criterion_text, sub_criteria_name)
        );`,
		`CREATE TABLE IF NOT EXISTS Bids (
            id INTEGER PRIMARY KEY AUTOINCREMENT, tender_id INTEGER NOT NULL, supplier_id INTEGER NOT NULL,
            submission_date TEXT DEFAULT CURRENT_TIMESTAMP, status TEXT DEFAULT 'submitted',
            financials_subtotal REAL NOT NULL, financials_vat REAL NOT NULL,
            financials_total REAL NOT NULL, financials_equivalentUSD REAL,
            validity_period_days INTEGER, cover_letter TEXT,
            FOREIGN KEY (tender_id) REFERENCES Tenders(id) ON DELETE CASCADE,
            FOREIGN KEY (supplier_id) REFERENCES Suppliers(id)
        );`,
		`CREATE TABLE IF NOT EXISTS BidItems (
            id INTEGER PRIMARY KEY AUTOINCREMENT, bid_id INTEGER NOT NULL, requisition_item_id INTEGER,
            itemName TEXT NOT NULL, description TEXT, quantity REAL NOT NULL, unitPrice REAL NOT NULL,
            totalPrice REAL NOT NULL, deliveryTerms_type TEXT, deliveryTerms_weeks INTEGER,
            deliveryTerms_location TEXT, warranty TEXT,
            FOREIGN KEY (bid_id) REFERENCES Bids(id) ON DELETE CASCADE,
            FOREIGN KEY (requisition_item_id) REFERENCES RequisitionItems(id) ON DELETE SET NULL
        );`,
		`CREATE TABLE IF NOT EXISTS BidEvaluationResults (
            id INTEGER PRIMARY KEY AUTOINCREMENT, bid_id INTEGER NOT NULL, criterion_id INTEGER NOT NULL,
            evaluator_id INTEGER, score REAL NOT NULL, comments TEXT,
            evaluation_date TEXT DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (bid_id) REFERENCES Bids(id) ON DELETE CASCADE,
            FOREIGN KEY (criterion_id) REFERENCES TenderEvaluationCriteria(id) ON DELETE CASCADE,
            FOREIGN KEY (evaluator_id) REFERENCES Users(id) ON DELETE SET NULL,
            UNIQUE (bid_id, criterion_id, evaluator_id)
        );`,
		`CREATE TABLE IF NOT EXISTS EvaluationPanels (
            id INTEGER PRIMARY KEY AUTOINCREMENT, tender_id INTEGER NOT NULL, panel_name TEXT,
            creation_date TEXT DEFAULT CURRENT_TIMESTAMP, status TEXT DEFAULT 'active',
            FOREIGN KEY (tender_id) REFERENCES Tenders(id) ON DELETE CASCADE
        );`,
		`CREATE TABLE IF NOT EXISTS UserEvaluationPanelMembership (
            user_id INTEGER NOT NULL, evaluation_panel_id INTEGER NOT NULL,
            role_in_panel TEXT, assigned_date TEXT DEFAULT CURRENT_TIMESTAMP,
            PRIMARY KEY (user_id, evaluation_panel_id),
            FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
            FOREIGN KEY (evaluation_panel_id) REFERENCES EvaluationPanels(id) ON DELETE CASCADE
        );`,
		`CREATE TABLE IF NOT EXISTS EvaluationPanelRecommendations (
            id INTEGER PRIMARY KEY AUTOINCREMENT, evaluation_panel_id INTEGER NOT NULL,
            recommended_bid_id INTEGER, recommendation_date TEXT DEFAULT CURRENT_TIMESTAMP,
            justification TEXT NOT NULL, overall_summary TEXT, status TEXT DEFAULT 'pending_approval',
            FOREIGN KEY (evaluation_panel_id) REFERENCES EvaluationPanels(id) ON DELETE CASCADE,
            FOREIGN KEY (recommended_bid_id) REFERENCES Bids(id) ON DELETE SET NULL
        );`,
		`CREATE TABLE IF NOT EXISTS EvaluationPanelSignatures (
            recommendation_id INTEGER NOT NULL, user_id INTEGER NOT NULL,
            digitalSignature TEXT, signed_at TEXT DEFAULT CURRENT_TIMESTAMP, comments TEXT,
            PRIMARY KEY (recommendation_id, user_id),
            FOREIGN KEY (recommendation_id) REFERENCES EvaluationPanelRecommendations(id) ON DELETE CASCADE,
            FOREIGN KEY (user_id) REFERENCES Users(id)
        );`,
		`CREATE TABLE IF NOT EXISTS PurchaseOrders (
            id INTEGER PRIMARY KEY AUTOINCREMENT, recommendation_id INTEGER UNIQUE, bid_id INTEGER UNIQUE,
            supplier_id INTEGER NOT NULL, po_number TEXT UNIQUE NOT NULL,
            order_date TEXT DEFAULT CURRENT_TIMESTAMP, status TEXT DEFAULT 'draft', 
            total_amount REAL NOT NULL, currency TEXT DEFAULT 'TZS', payment_terms TEXT, delivery_address TEXT,
            data_plant TEXT, data_purchasingGroup TEXT, data_wbsElement TEXT, created_by_user_id INTEGER,
            FOREIGN KEY (recommendation_id) REFERENCES EvaluationPanelRecommendations(id) ON DELETE SET NULL,
            FOREIGN KEY (bid_id) REFERENCES Bids(id) ON DELETE SET NULL,
            FOREIGN KEY (supplier_id) REFERENCES Suppliers(id),
            FOREIGN KEY (created_by_user_id) REFERENCES Users(id) ON DELETE SET NULL
        );`,
		`CREATE TABLE IF NOT EXISTS PurchaseOrderItems (
            id INTEGER PRIMARY KEY AUTOINCREMENT, purchase_order_id INTEGER NOT NULL, bid_item_id INTEGER,
            requisition_item_id INTEGER, description TEXT NOT NULL, quantity REAL NOT NULL,
            unit_price REAL NOT NULL, total_price REAL NOT NULL, amr_id INTEGER,
            ancillaryCostsAllocated REAL DEFAULT 0, capitalizedValue REAL, delivery_status TEXT DEFAULT 'pending',
            FOREIGN KEY (purchase_order_id) REFERENCES PurchaseOrders(id) ON DELETE CASCADE,
            FOREIGN KEY (bid_item_id) REFERENCES BidItems(id) ON DELETE SET NULL,
            FOREIGN KEY (requisition_item_id) REFERENCES RequisitionItems(id) ON DELETE SET NULL,
            FOREIGN KEY (amr_id) REFERENCES Assets(id) ON DELETE SET NULL
        );`,
		`CREATE TABLE IF NOT EXISTS Deliveries (
            id INTEGER PRIMARY KEY AUTOINCREMENT, purchase_order_id INTEGER NOT NULL,
            purchase_order_item_id INTEGER, grnNumber TEXT UNIQUE NOT NULL,
            delivery_date TEXT DEFAULT CURRENT_TIMESTAMP, received_by_user_id INTEGER, quantity_received REAL,
            asset_id_received INTEGER, assetCommissioningDate TEXT, assetDeliveryDate TEXT,
            complianceStatus TEXT, capitalizationStatus TEXT, notes TEXT,
            FOREIGN KEY (purchase_order_id) REFERENCES PurchaseOrders(id) ON DELETE RESTRICT,
            FOREIGN KEY (purchase_order_item_id) REFERENCES PurchaseOrderItems(id) ON DELETE SET NULL,
            FOREIGN KEY (received_by_user_id) REFERENCES Users(id) ON DELETE SET NULL,
            FOREIGN KEY (asset_id_received) REFERENCES Assets(id) ON DELETE SET NULL
        );`,
		`CREATE TABLE IF NOT EXISTS Invoices (
            id INTEGER PRIMARY KEY AUTOINCREMENT, delivery_id INTEGER UNIQUE, purchase_order_id INTEGER,
            supplier_id INTEGER NOT NULL, invoice_number TEXT UNIQUE NOT NULL, invoice_date TEXT NOT NULL,
            due_date TEXT, subtotal REAL NOT NULL, tax_amount REAL, total_amount REAL NOT NULL,
            status TEXT DEFAULT 'pending', payment_date TEXT, payment_reference TEXT,
            FOREIGN KEY (delivery_id) REFERENCES Deliveries(id) ON DELETE SET NULL,
            FOREIGN KEY (purchase_order_id) REFERENCES PurchaseOrders(id) ON DELETE RESTRICT,
            FOREIGN KEY (supplier_id) REFERENCES Suppliers(id)
        );`,
		`CREATE TABLE IF NOT EXISTS Attachments (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            entity_type TEXT NOT NULL, entity_id INTEGER NOT NULL,
            file_name TEXT NOT NULL, file_path TEXT, description TEXT,
            uploaded_by_user_id INTEGER, upload_date TEXT DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (uploaded_by_user_id) REFERENCES Users(id) ON DELETE SET NULL
        );`,
	}

	for _, tableSQL := range tables {
		_, err = db.Exec(tableSQL)
		if err != nil {
			log.Printf("Error creating table: %v\nSQL: %s\n", err, tableSQL)
			// Decide if you want to stop on first error or continue
			// log.Fatalf("Stopping due to error creating table: %v", err)
		}
	}

	log.Println("Database tables created/verified successfully.")
}
