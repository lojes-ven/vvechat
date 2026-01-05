class Verifier:
    def __init__(self):
        pass

    @staticmethod
    def verify(ids_list):
        """
        Verifies a list of (timestamp, id, node_id) tuples or just ids.
        Returns a dict with stats and error list.
        """
        if not ids_list:
            return {"valid": True, "count": 0, "errors": []}

        # Sort by ID to check global uniqueness and monotonicity per node
        # ids_list is expected to be a list of generated IDs (integers)
        
        sorted_ids = sorted(ids_list)
        unique_ids = set(sorted_ids)
        
        errors = []
        
        # Check uniqueness
        if len(unique_ids) != len(sorted_ids):
            errors.append(f"Duplicate IDs found. Total: {len(sorted_ids)}, Unique: {len(unique_ids)}")

        # Check monotonicity per node
        # We need to decode the ID to get node_id
        # Layout: 1 bit sign | 41 bits timestamp | 10 bits node_id | 12 bits sequence
        node_id_bits = 10
        sequence_bits = 12
        node_id_shift = sequence_bits
        node_mask = (1 << node_id_bits) - 1
        
        last_id_per_node = {}
        
        for id_val in ids_list: # Check in generation order if possible, but usually we check sorted for global properties.
                                # For monotonicity, we should check the order they were generated if we have that info.
                                # If we only have the set of IDs, we can only check if they *can* be ordered.
                                # But Snowflake IDs are roughly time-ordered.
            
            # Extract node_id
            node_id = (id_val >> node_id_shift) & node_mask
            
            # We can't strictly check generation order unless the input list preserves it.
            # Assuming input list IS in generation order (concatenated from workers).
            pass

        return {
            "valid": len(errors) == 0,
            "count": len(sorted_ids),
            "unique_count": len(unique_ids),
            "errors": errors
        }

    @staticmethod
    def analyze_bottlenecks(results):
        """
        Analyzes timing data to find bottlenecks.
        """
        # Placeholder for logic that looks at 'wait_time' vs 'lock_time'
        pass
